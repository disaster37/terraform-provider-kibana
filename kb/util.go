package kb

import "encoding/json"

// optionalInterfaceJson permit to convert string as json object
func optionalInterfaceJson(input string) interface{} {
	if input == "" || input == "{}" {
		return nil
	} else {
		return json.RawMessage(input)
	}
}

// convertArrayInterfaceToArrayString permit to convert an array of interface to an array of string
func convertArrayInterfaceToArrayString(raws []interface{}) []string {
	data := make([]string, len(raws))
	for i, raw := range raws {
		data[i] = raw.(string)
	}

	return data
}

func convertMapInterfaceToMapString(raws map[string]interface{}) map[string]string {
	data := make(map[string]string)
	for k, v := range raws {
		data[k] = v.(string)
	}

	return data
}
