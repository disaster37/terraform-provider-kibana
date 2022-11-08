package kb

import (
	"encoding/json"
	"fmt"
	"strings"

	eshandler "github.com/disaster37/es-handler/v8"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

// suppressEquivalentJSON permit to compare state store as JSON string
func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {

	var err error
	oldObj := map[string]any{}
	newObj := map[string]any{}

	if old == "" {
		old = "{}"
	}
	if new == "" {
		new = "{}"
	}

	if err = json.Unmarshal([]byte(old), &oldObj); err != nil {
		fmt.Printf("[ERR] Error when converting current Json: %s\ndata: %s", err.Error(), old)
		log.Errorf("Error when converting current Json: %s\ndata: %s", err.Error(), old)
	}
	if err = json.Unmarshal([]byte(new), &newObj); err != nil {
		fmt.Printf("[ERR] Error when converting current Json: %s\ndata: %s", err.Error(), new)
		log.Errorf("Error when converting current Json: %s\ndata: %s", err.Error(), new)
	}

	diff, err := eshandler.StandardDiff(oldObj, newObj, logEntry, nil)
	if err != nil {
		fmt.Printf("[ERR] Error when diff JSON: %s", err.Error())
		log.Errorf("Error when diff Json: %s", err.Error())
		return false
	}

	return diff == ""
}

func suppressEquivalentJSONWithExclude(k, old, new string, d *schema.ResourceData, exclude map[string]any) bool {

	var err error
	oldObj := map[string]any{}
	newObj := map[string]any{}

	if old == "" {
		old = "{}"
	}
	if new == "" {
		new = "{}"
	}

	if err = json.Unmarshal([]byte(old), &oldObj); err != nil {
		fmt.Printf("[ERR] Error when converting current Json: %s\ndata: %s", err.Error(), old)
		log.Errorf("Error when converting current Json: %s\ndata: %s", err.Error(), old)
	}
	if err = json.Unmarshal([]byte(new), &newObj); err != nil {
		fmt.Printf("[ERR] Error when converting current Json: %s\ndata: %s", err.Error(), new)
		log.Errorf("Error when converting current Json: %s\ndata: %s", err.Error(), new)
	}

	diff, err := eshandler.StandardDiff(oldObj, newObj, logEntry, exclude)
	if err != nil {
		fmt.Printf("[ERR] Error when diff JSON: %s", err.Error())
		log.Errorf("Error when diff Json: %s", err.Error())
		return false
	}

	return diff == ""
}

// Split NDJson by keeping only not emty lines
func splitNDJSON(val string) []string {
	slices := strings.Split(val, "\n")
	result := []string{}

	for i := range slices {
		if len(slices[i]) > 0 {
			result = append(result, slices[i])
		}
	}

	return result
}

// suppressEquivalentNDJSON permit to compare ndjson string
func suppressEquivalentNDJSON(k, old, new string, d *schema.ResourceData) bool {

	var err error
	excludeFields := map[string]any{
		"version":              nil,
		"updated_at":           nil,
		"coreMigrationVersion": nil,
		"migrationVersion":     nil,
		"references":           nil,
		"sort":                 nil,
	}

	// NDJSON mean sthat each line correspond to JSON struct
	oldSlice := splitNDJSON(old)
	newSlice := splitNDJSON(new)
	oldObjSlice := make([]map[string]any, len(oldSlice))
	newObjSlice := make([]map[string]any, len(newSlice))
	if len(oldSlice) != len(newSlice) {
		return false
	}

	// Convert each line to map of string to compare the same object id
	for i, oldJSON := range oldSlice {
		res := map[string]any{}
		if oldJSON != "" {
			if err = json.Unmarshal([]byte(oldJSON), &res); err != nil {
				fmt.Printf("[ERR] Error when unmarshal old Json: %s\ndata: %s", err.Error(), oldJSON)
				log.Errorf("Error when unmarshal old Json: %s\ndata: %s", err.Error(), oldJSON)
				return false
			}
		}
		oldObjSlice[i] = res
	}

	for i, newJSON := range newSlice {
		res := map[string]any{}
		if newJSON != "" {
			if err = json.Unmarshal([]byte(newJSON), &res); err != nil {
				fmt.Printf("[ERR] Error when unmarshal nes Json: %s\ndata: %s", err.Error(), newJSON)
				log.Errorf("Error when unmarshal new Json: %s\ndata: %s", err.Error(), newJSON)
				return false

			}
		}
		newObjSlice[i] = res
	}

	// Compare json obj
	for i, oldItem := range oldObjSlice {
		isFound := false
		if oldItem["id"] == "" {
			return false
		}
		for j, newItem := range newObjSlice {
			if newItem["id"] == "" {
				return false
			}

			// Compare same items
			if oldItem["id"] == newItem["id"] {
				if !suppressEquivalentJSONWithExclude(k, oldSlice[i], newSlice[j], d, excludeFields) {
					return false
				}
				isFound = true
				break
			}
		}

		if !isFound {
			return false
		}
	}

	return true
}
