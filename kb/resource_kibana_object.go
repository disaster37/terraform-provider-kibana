// Manage the object in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/saved-objects-api.html
// Supported version:
//  - v7
package kb

import (
	"fmt"

	kibana "github.com/disaster37/go-kibana-rest/v7"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle kibana save object
func resourceKibanaObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaObjectCreate,
		Read:   resourceKibanaObjectRead,
		Update: resourceKibanaObjectUpdate,
		Delete: resourceKibanaObjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"space": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
			"data": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentJson,
			},
			"export_types": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"export_objects": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"deep_reference": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

// Import objects in Kibana
func resourceKibanaObjectCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	err := importObject(d, meta)
	if err != nil {
		return err
	}

	d.SetId(name)

	log.Infof("Imported objects %s successfully", name)

	return resourceKibanaObjectRead(d, meta)
}

// Export objects in Kibana
func resourceKibanaObjectRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	exportTypes := convertArrayInterfaceToArrayString(d.Get("export_types").(*schema.Set).List())
	exportObjects := buildExportObjects(d.Get("export_objects").(*schema.Set).List())
	deepReference := d.Get("deep_reference").(bool)
	space := d.Get("space").(string)

	log.Debugf("Object id:  %s", id)
	log.Debugf("Export types: %+v", exportTypes)
	log.Debugf("Export Objects: %+v", exportObjects)
	log.Debugf("Space: %s", space)

	client := meta.(*kibana.Client)

	data, err := client.API.KibanaSavedObject.Export(exportTypes, exportObjects, deepReference, space)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		fmt.Printf("[WARN] Export object %s not found - removing from state", id)
		log.Warnf("Export object %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	log.Debugf("Export object %s successfully:\n%+v", id, data)

	d.Set("name", id)
	d.Set("data", data)
	d.Set("space", space)
	d.Set("export_types", exportTypes)
	d.Set("export_objects", exportObjects)

	log.Infof("Export object %s successfully", id)

	return nil
}

// Update existing object in Kibana
func resourceKibanaObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := importObject(d, meta)
	if err != nil {
		return err
	}

	log.Infof("Updated object %s successfully", id)

	return resourceKibanaObjectRead(d, meta)
}

// Delete object in Kibana is not supported
// It just remove object from state
func resourceKibanaObjectDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	log.Infof("Delete object in not supported")
	fmt.Printf("[INFO] Delete object in not supported - just removing from state")
	return nil

}

// Build list of object to export
func buildExportObjects(raws []interface{}) []map[string]string {

	results := make([]map[string]string, len(raws))

	for i, raw := range raws {
		m := raw.(map[string]interface{})
		object := map[string]string{}
		object["type"] = m["type"].(string)
		object["id"] = m["id"].(string)
		results[i] = object
	}

	return results

}

// Import objects in Kibana
func importObject(d *schema.ResourceData, meta interface{}) error {
	data := d.Get("data").(string)
	space := d.Get("space").(string)

	log.Debugf("Data: %s", data)

	var (
		importedData map[string]interface{}
		err          error
	)

	client := meta.(*kibana.Client)

	importedData, err = client.API.KibanaSavedObject.Import([]byte(data), true, space)
	if err != nil {
		return err
	}

	log.Debugf("Imported object: %+v", importedData)

	return nil
}
