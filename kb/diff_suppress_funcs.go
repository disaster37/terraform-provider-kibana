package kb

import (
	"strings"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/diff"
	ucfgjson "github.com/elastic/go-ucfg/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

// suppressEquivalentJSON permit to compare json string
func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {
	confOld, err := ucfgjson.NewConfig([]byte(old), ucfg.PathSep("."))
	if err != nil {
		log.Errorf("Error on suppressEquivalentJSON: %s", err.Error())
		return false
	}
	confNew, err := ucfgjson.NewConfig([]byte(new), ucfg.PathSep("."))
	if err != nil {
		log.Errorf("Error on suppressEquivalentJSON: %s", err.Error())
		return false
	}

	currentDiff := diff.CompareConfigs(confOld, confNew)
	log.Debugf("Diff\n: %s", currentDiff.GoStringer())

	return !currentDiff.HasChanged()
}

// suppressEquivalentNDJSON permit to compare ndjson string
func suppressEquivalentNDJSON(k, old, new string, d *schema.ResourceData) bool {

	// NDJSON mean sthat each line correspond to JSON struct
	oldSlice := strings.Split(old, "\n")
	newSlice := strings.Split(new, "\n")
	oldObjSlice := make([]*ucfg.Config, len(oldSlice))
	newObjSlice := make([]*ucfg.Config, len(newSlice))
	if len(oldSlice) != len(newSlice) {
		return false
	}

	// Convert string line to JSON
	for i, oldJSON := range oldSlice {
		config, err := ucfgjson.NewConfig([]byte(oldJSON), ucfg.PathSep("."))
		if err != nil {
			log.Errorf("Error on suppressEquivalentNDJSON: %s", err.Error())
			return false
		}
		config.Remove("version", -1)
		config.Remove("updated_at", -1)

		oldObjSlice[i] = config
	}
	for i, newJSON := range newSlice {
		config, err := ucfgjson.NewConfig([]byte(newJSON), ucfg.PathSep("."))
		if err != nil {
			log.Errorf("Error on suppressEquivalentNDJSON: %s", err.Error())
			return false
		}
		config.Remove("version", -1)
		config.Remove("updated_at", -1)

		newObjSlice[i] = config
	}

	// Compare json obj
	for _, oldConfig := range oldObjSlice {
		isFound := false
		oldId, err := oldConfig.String("id", -1)
		if err != nil {
			log.Errorf("Error on suppressEquivalentNDJSON: %s", err.Error())
			return false
		}
		for _, newConfig := range newObjSlice {
			newId, err := newConfig.String("id", -1)
			if err != nil {
				log.Errorf("Error on suppressEquivalentNDJSON: %s", err.Error())
				return false
			}
			if oldId == newId {
				currentDiff := diff.CompareConfigs(oldConfig, newConfig)
				log.Debugf("Diff\n: %s", currentDiff.GoStringer())

				if currentDiff.HasChanged() {
					return false
				}
				isFound = true
				break
			}
		}

		if isFound == false {
			return false
		}
	}

	/*
		// NDJSON mean sthat each line correspond to JSON struct
		oldSlice := strings.Split(old, "\n")
		newSlice := strings.Split(new, "\n")
		oldObjSlice := make([]map[string]interface{}, len(oldSlice))
		newObjSlice := make([]map[string]interface{}, len(newSlice))
		if len(oldSlice) != len(newSlice) {
			return false
		}

		// Convert string line to JSON
		for i, oldJSON := range oldSlice {
			jsonObj := make(map[string]interface{})
			if err := json.Unmarshal([]byte(oldJSON), &jsonObj); err != nil {
				return false
			}

			delete(jsonObj, "version")
			delete(jsonObj, "updated_at")

			oldObjSlice[i] = jsonObj
		}
		for i, newJSON := range newSlice {
			jsonObj := make(map[string]interface{})
			if err := json.Unmarshal([]byte(newJSON), &jsonObj); err != nil {
				return false
			}
			delete(jsonObj, "version")
			delete(jsonObj, "updated_at")

			newObjSlice[i] = jsonObj
		}

		// Compare json obj
		for _, oldJSON := range oldObjSlice {
			isFound := false
			for _, newJSON := range newObjSlice {
				if oldJSON["id"].(string) == newJSON["id"].(string) {
					if reflect.DeepEqual(oldJSON, newJSON) == false {
						return false
					}
					isFound = true
					break
				}
			}

			if isFound == false {
				return false
			}
		}
	*/

	return true

}
