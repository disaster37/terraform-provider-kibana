package kb

import (
	"crypto/tls"
	"net/http"
	"net/url"

	kibana7 "github.com/disaster37/go-kibana-rest"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_URL", nil),
				Description: "Kibana URL",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_USERNAME", nil),
				Description: "Username to use to connect to Kibana using basic auth",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("KIBANA_PASSWORD", nil),
				Description: "Password to use to connect to Kibana using basic auth",
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "A Custom CA certificate",
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable SSL verification of API calls",
			},
		},

		ResourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	var (
		relevantClient interface{}
		data           map[string]interface{}
	)

	URL := d.Get("url").(string)
	insecure := d.Get("insecure").(bool)
	cacertFile := d.Get("cacert_file").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}
	// Checks is valid URL
	if _, err := url.Parse(URL); err != nil {
		return nil, err
	}

	// Intialise connexion
	cfg := kibana7.Config{
		Address: URL,
	}
	if username != "" && password != "" {
		cfg.Username = username
		cfg.Password = password
	}
	if insecure == true {
		cfg.DisableVerifySSL = true
	}
	// If a cacertFile has been specified, use that for cert validation
	/*
		if cacertFile != "" {
			caCert, _, _ := pathorcontents.Read(cacertFile)

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM([]byte(caCert))
			transport.TLSClientConfig.RootCAs = caCertPool
		}
		cfg.Transport = transport
	*/
	client, err := kibana7.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Test connexion and check elastic version to use the right Version
	kibanaStatus, err := client.API.KibanaStatus.Get()
	if err != nil {
		return nil, err
	}

	if kibanaStatus == nil {
		return nil, errors.New("Status is empty, somethink wrong with Kibana ?")
	}

	version := kibanaStatus["version"].(map[string]interface{})["number"].(string)
	log.Debugf("Server: %s", version)

	if version < "8.0.0" && version >= "7.0.0" {
		log.Printf("[INFO] Using Kibana 7")
		relevantClient = client
	} else {
		return nil, errors.New("Kibana is older than 7.0.0!")
	}

	return relevantClient, nil
}
