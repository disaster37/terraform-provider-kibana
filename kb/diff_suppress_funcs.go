package kb

import (
	"encoding/json"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// suppressEquivalentJSON permit to compare json string
func suppressEquivalentJSON(k, old, new string, d *schema.ResourceData) bool {
	var oldObj, newObj interface{}
	if err := json.Unmarshal([]byte(old), &oldObj); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newObj); err != nil {
		return false
	}
	return reflect.DeepEqual(oldObj, newObj)
}
