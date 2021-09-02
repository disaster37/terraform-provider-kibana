TEST?=./...
PKG_NAME=kb
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
KIBANA_URL ?= http://127.0.0.1:5601
KIBANA_USERNAME ?= elastic
KIBANA_PASSWORD ?= changeme
ELASTICSEARCH_URLS ?= http://127.0.0.1:9200
ELASTICSEARCH_USERNAME ?= elastic
ELASTICSEARCH_PASSWORD ?= changeme

default: build

build: fmt fmtcheck
	go install

local-build:
	mkdir -p registry/registry.terraform.io/disaster37/kibana/1.0.0/linux_amd64
	go build -o registry/registry.terraform.io/disaster37/kibana/1.0.0/linux_amd64/terraform-provider-kibana

gen:
	rm -f aws/internal/keyvaluetags/*_gen.go
	go generate ./...

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmt fmtcheck
	KIBANA_URL=${KIBANA_URL} KIBANA_USERNAME=${KIBANA_USERNAME} KIBANA_PASSWORD=${KIBANA_PASSWORD} TF_ACC=1 TF_LOG_PROVIDER=DEBUG go test $(TEST) -v -count 1 -parallel 1 -race -coverprofile=coverage.txt -covermode=atomic $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

websitefmtcheck:
	@sh -c "'$(CURDIR)/scripts/websitefmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run --no-config --deadline 5m --disable-all --enable staticcheck --exclude SA1019 --max-issues-per-linter 0 --max-same-issues 0 ./$(PKG_NAME)/...
	@golangci-lint run ./$(PKG_NAME)/...
	@tfproviderlint \
		-c 1 \
		-AT001 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		./$(PKG_NAME)

tools:
	GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

start-pods: clean-pods
	kubectl run elasticsearch --image docker.elastic.co/elasticsearch/elasticsearch:7.12.1 --port "9200" --expose --env "cluster.name=test" --env "discovery.type=single-node" --env "ELASTIC_PASSWORD=changeme" --env "xpack.security.enabled=true" --env "ES_JAVA_OPTS=-Xms512m -Xmx512m" --env "path.repo=/tmp" --limits "cpu=500m,memory=1024Mi"
	kubectl run kibana --image docker.elastic.co/kibana/kibana:7.12.1 --expose --port "5601" --env "ELASTICSEARCH_HOSTS=http://elasticsearch:9200" --env "ELASTICSEARCH_USERNAME=elastic" --env "ELASTICSEARCH_PASSWORD=changeme" --limits "cpu=500m,memory=512Mi"

clean-pods:
	kubectl delete --ignore-not-found pod/kibana
	kubectl delete --ignore-not-found service/kibana
	kubectl delete --ignore-not-found pod/elasticsearch
	kubectl delete --ignore-not-found service/elasticsearch

trial-license:
	curl -XPOST -u ${ELASTICSEARCH_USERNAME}:${ELASTICSEARCH_PASSWORD} ${ELASTICSEARCH_URLS}/_license/start_trial?acknowledge=true

.PHONY: build gen sweep test testacc fmt fmtcheck lint tools test-compile website website-lint website-test start-pods clean-pods local-build trial-license