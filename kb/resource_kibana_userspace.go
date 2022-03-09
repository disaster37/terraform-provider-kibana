// Manage the user space in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/spaces-api.html
// Supported version:
//  - v7

package kb

import (
	"fmt"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	kbapi "github.com/disaster37/go-kibana-rest/v7/kbapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle user space in Kibana
func resourceKibanaUserSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaUserSpaceCreate,
		Read:   resourceKibanaUserSpaceRead,
		Update: resourceKibanaUserSpaceUpdate,
		Delete: resourceKibanaUserSpaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled_features": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"initials": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// Create new user space in Kibana
func resourceKibanaUserSpaceCreate(d *schema.ResourceData, meta interface{}) error {
	id := d.Get("uid").(string)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	disabledFeatures := convertArrayInterfaceToArrayString(d.Get("disabled_features").(*schema.Set).List())
	initials := d.Get("initials").(string)
	color := d.Get("color").(string)

	client := meta.(*kibana.Client)

	userSpace := &kbapi.KibanaSpace{
		ID:               id,
		Name:             name,
		Description:      description,
		DisabledFeatures: disabledFeatures,
		Initials:         initials,
		Color:            color,
	}

	userSpace, err := client.API.KibanaSpaces.Create(userSpace)
	if err != nil {
		return err
	}

	d.SetId(id)

	log.Infof("Created user space %s (%s) successfully", id, name)
	fmt.Printf("[INFO] Created user space %s (%s) successfully", id, name)

	return resourceKibanaUserSpaceRead(d, meta)
}

// Read existing user space in Kibana
func resourceKibanaUserSpaceRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	log.Debugf("User space id:  %s", id)

	client := meta.(*kibana.Client)

	userSpace, err := client.API.KibanaSpaces.Get(id)
	if err != nil {
		return err
	}

	if userSpace == nil {
		log.Warnf("User space %s not found - removing from state", id)
		fmt.Printf("[WARN] User space %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	log.Debugf("Get user space %s successfully:\n%s", id, userSpace)

	d.Set("uid", id)

	d.Set("name", userSpace.Name)
	d.Set("description", userSpace.Description)
	d.Set("disabled_features", userSpace.DisabledFeatures)
	d.Set("initials", userSpace.Initials)
	d.Set("color", userSpace.Color)

	log.Infof("Read user space %s successfully", id)
	fmt.Printf("[INFO] Read user space %s successfully", id)

	return nil
}

// Update existing user space in Elasticsearch
func resourceKibanaUserSpaceUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	disabledFeatures := convertArrayInterfaceToArrayString(d.Get("disabled_features").(*schema.Set).List())
	initials := d.Get("initials").(string)
	color := d.Get("color").(string)

	client := meta.(*kibana.Client)
	userSpace := &kbapi.KibanaSpace{
		ID:               id,
		Name:             name,
		Description:      description,
		DisabledFeatures: disabledFeatures,
		Initials:         initials,
		Color:            color,
	}

	userSpace, err := client.API.KibanaSpaces.Update(userSpace)
	if err != nil {
		return err
	}

	log.Infof("Updated user space %s successfully", id)
	fmt.Printf("[INFO] Updated user space %s successfully", id)

	return resourceKibanaUserSpaceRead(d, meta)
}

// Delete existing user space in Kibana
func resourceKibanaUserSpaceDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	log.Debugf("User space id: %s", id)

	client := meta.(*kibana.Client)

	err := client.API.KibanaSpaces.Delete(id)
	if err != nil {
		if err.(kbapi.APIError).Code == 404 {
			log.Warnf("User space %s not found - removing from state", id)
			fmt.Printf("[WARN] User space %s not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return err

	}

	d.SetId("")

	log.Infof("Deleted user space %s successfully", id)
	fmt.Printf("[INFO] Deleted user space %s successfully", id)
	return nil

}
