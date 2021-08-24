// Copy Kibana object from space to another spaces
// API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api-copy-saved-objects.html
// Supported version:
//  - v7

package kb

import (
	"log"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	"github.com/disaster37/go-kibana-rest/v7/kbapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Resource specification to handle kibana save object
func resourceKibanaCopyObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaCopyObjectCreate,
		Read:   resourceKibanaCopyObjectRead,
		Update: resourceKibanaCopyObjectUpdate,
		Delete: resourceKibanaCopyObjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_space": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"target_spaces": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"object": {
				Type:     schema.TypeSet,
				Required: true,
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
			"include_reference": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"force_update": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Copy objects in Kibana
func resourceKibanaCopyObjectCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)

	err := copyObject(d, meta)
	if err != nil {
		return err
	}

	d.SetId(name)

	log.Printf("[INFO] Copy objects %s successfully", name)

	return resourceKibanaCopyObjectRead(d, meta)
}

// Read object on kibana
func resourceKibanaCopyObjectRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	sourceSpace := d.Get("source_space").(string)
	targetSpaces := convertArrayInterfaceToArrayString(d.Get("target_spaces").(*schema.Set).List())
	objects := buildCopyObjects(d.Get("object").(*schema.Set).List())
	includeReference := d.Get("include_reference").(bool)
	overwrite := d.Get("overwrite").(bool)
	forceUpdate := d.Get("force_update").(bool)

	log.Printf("[DEBUG] Resource id:  %s", id)
	log.Printf("[DEBUG] Source space: %s", sourceSpace)
	log.Printf("[DEBUG] Target spaces: %+v", targetSpaces)
	log.Printf("[DEBUG] Objects: %+v", objects)
	log.Printf("[DEBUG] Include reference: %t", includeReference)
	log.Printf("[DEBUG] Overwrite: %t", overwrite)
	log.Printf("[DEBUG] force_update: %t", forceUpdate)

	// @ TODO
	// A good when is to check if already exported object is the same that original space
	// To avoid this hard code, we juste use force_update and overwrite field
	// It make same result but in bad way on terraform spirit

	d.Set("name", id)
	d.Set("source_space", sourceSpace)
	d.Set("target_spaces", targetSpaces)
	d.Set("object", objects)
	d.Set("include_reference", includeReference)
	d.Set("overwrite", overwrite)
	d.Set("force_update", false)

	log.Printf("[INFO] Read resource %s successfully", id)

	return nil
}

// Update existing object in Kibana
func resourceKibanaCopyObjectUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := copyObject(d, meta)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updated resource %s successfully", id)

	return resourceKibanaCopyObjectRead(d, meta)
}

// Delete object in Kibana is not supported
// It just remove object from state
func resourceKibanaCopyObjectDelete(d *schema.ResourceData, meta interface{}) error {

	d.SetId("")

	log.Printf("[INFO] Delete object in not supported - just removing from state")
	return nil

}

// Build list of object to export
func buildCopyObjects(raws []interface{}) []map[string]string {

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

// Copy objects in Kibana
func copyObject(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	sourceSpace := d.Get("source_space").(string)
	targetSpaces := convertArrayInterfaceToArrayString(d.Get("target_spaces").(*schema.Set).List())
	objects := buildCopyObjects(d.Get("object").(*schema.Set).List())
	includeReference := d.Get("include_reference").(bool)
	overwrite := d.Get("overwrite").(bool)

	log.Printf("[DEBUG] Source space: %s", sourceSpace)
	log.Printf("[DEBUG] Target spaces: %+v", targetSpaces)
	log.Printf("[DEBUG] Objects: %+v", objects)
	log.Printf("[DEBUG] Include reference: %t", includeReference)
	log.Printf("[DEBUG] Overwrite: %t", overwrite)

	client := meta.(*kibana.Client)

	objectsParameter := make([]kbapi.KibanaSpaceObjectParameter, 0, 1)
	for _, object := range objects {
		objectsParameter = append(objectsParameter, kbapi.KibanaSpaceObjectParameter{
			ID:   object["id"],
			Type: object["type"],
		})
	}

	parameter := &kbapi.KibanaSpaceCopySavedObjectParameter{
		Spaces:            targetSpaces,
		Objects:           objectsParameter,
		IncludeReferences: includeReference,
		Overwrite:         overwrite,
	}

	err := client.API.KibanaSpaces.CopySavedObjects(parameter, sourceSpace)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Copy object for resource successfully: %s", name)

	return nil
}
