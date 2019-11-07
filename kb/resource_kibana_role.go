// Manage the role in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/role-management-api.html
// Supported version:
//  - v7
package kb

import (
	"fmt"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	kbapi "github.com/disaster37/go-kibana-rest/v7/kbapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle role in Kibana
func resourceKibanaRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaRoleCreate,
		Read:   resourceKibanaRoleRead,
		Update: resourceKibanaRoleUpdate,
		Delete: resourceKibanaRoleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"elasticsearch": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"indices": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"names": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"privileges": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"query": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "{}",
										DiffSuppressFunc: suppressEquivalentJson,
									},
									"field_security": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "{}",
										DiffSuppressFunc: suppressEquivalentJson,
									},
								},
							},
						},
						"cluster": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"run_as": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"kibana": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"spaces": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"features": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"permissions": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "{}",
				DiffSuppressFunc: suppressEquivalentJson,
			},
		},
	}
}

// Create new role in Kibana
func resourceKibanaRoleCreate(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)

	err := createRole(d, meta)
	if err != nil {
		return err
	}

	d.SetId(name)

	log.Infof("Created role %s successfully", name)

	return resourceKibanaRoleRead(d, meta)
}

// Read existing role in Kibana
func resourceKibanaRoleRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	log.Debugf("Role id:  %s", id)

	client := meta.(*kibana.Client)

	role, err := client.API.KibanaRoleManagement.Get(id)
	if err != nil {
		return err
	}

	if role == nil {
		fmt.Printf("[WARN] Role %s not found - removing from state", id)
		log.Warnf("Role %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	log.Debugf("Get role %s successfully:\n%s", id, role)

	d.Set("name", id)
	d.Set("elasticsearch", []kbapi.KibanaRoleElasticsearch{*role.Elasticsearch})
	d.Set("kibana", role.Kibana)
	d.Set("metadata", role.Metadata)

	log.Infof("Read role %s successfully", id)

	return nil
}

// Update existing role in Elasticsearch
func resourceKibanaRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	err := createRole(d, meta)
	if err != nil {
		return err
	}

	log.Infof("Updated role %s successfully", id)

	return resourceKibanaRoleRead(d, meta)
}

// Delete existing role in Elasticsearch
func resourceKibanaRoleDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	log.Debugf("Role id: %s", id)

	client := meta.(*kibana.Client)

	err := client.API.KibanaRoleManagement.Delete(id)
	if err != nil {
		if err.(kbapi.APIError).Code == 404 {
			fmt.Printf("[WARN] Role %s not found - removing from state", id)
			log.Warnf("Role %s not found - removing from state", id)
			d.SetId("")
			return nil
		} else {
			return err
		}
	}

	d.SetId("")

	log.Infof("Deleted role %s successfully", id)
	return nil

}

// createRole permit to create or update role in Kibana
func createRole(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	metadataTemp := optionalInterfaceJson(d.Get("metadata").(string))
	roleElasticsearch := buildRolesElasticsearch(d.Get("elasticsearch").(*schema.Set).List())
	roleKibana := buildRolesKibana(d.Get("kibana").(*schema.Set).List())

	client := meta.(*kibana.Client)

	var metadata map[string]interface{}
	if metadataTemp != nil {
		metadata = metadataTemp.(map[string]interface{})
	} else {
		metadata = nil
	}
	role := &kbapi.KibanaRole{
		Name:          name,
		Elasticsearch: roleElasticsearch,
		Kibana:        roleKibana,
		Metadata:      metadata,
	}

	role, err := client.API.KibanaRoleManagement.CreateOrUpdate(role)
	if err != nil {
		return err
	}

	return nil
}

// buildRolesElasticsearch permit to construct kibanaRoleElasticsearch object
func buildRolesElasticsearch(raws []interface{}) *kbapi.KibanaRoleElasticsearch {
	if len(raws) == 0 {
		return nil
	}

	// We check only the first, we case use multiple KibanaRoleElasticsearch
	raw := raws[0].(map[string]interface{})

	kibanaRoleElasticsearch := &kbapi.KibanaRoleElasticsearch{}

	if _, ok := raw["run_as"]; ok {
		kibanaRoleElasticsearch.RunAs = convertArrayInterfaceToArrayString(raw["run_as"].(*schema.Set).List())
	}
	if _, ok := raw["cluster"]; ok {
		kibanaRoleElasticsearch.Cluster = convertArrayInterfaceToArrayString(raw["cluster"].(*schema.Set).List())
	}
	if _, ok := raw["indices"]; ok {
		kibanaRoleElasticsearch.Indices = buildKibanaRoleElasticsearchIndice(raw["indices"].(*schema.Set).List())
	}

	return kibanaRoleElasticsearch

}

// buildKibanaRoleElasticsearchIndice permit to build list of KibanaRoleElasticsearchIndice
func buildKibanaRoleElasticsearchIndice(raws []interface{}) []kbapi.KibanaRoleElasticsearchIndice {
	kibanaRoleElasticsearchIndices := make([]kbapi.KibanaRoleElasticsearchIndice, len(raws))
	for i, raw := range raws {
		m := raw.(map[string]interface{})
		fieldSecurityTemp := optionalInterfaceJson(m["field_security"].(string))
		var fieldSecurity map[string]interface{}
		if fieldSecurityTemp != nil {
			fieldSecurity = fieldSecurityTemp.(map[string]interface{})
		} else {
			fieldSecurity = nil
		}
		kibanaRoleElasticsearchIndice := kbapi.KibanaRoleElasticsearchIndice{
			Names:         convertArrayInterfaceToArrayString(m["names"].(*schema.Set).List()),
			Privileges:    convertArrayInterfaceToArrayString(m["privileges"].(*schema.Set).List()),
			Query:         optionalInterfaceJson(m["query"].(string)),
			FieldSecurity: fieldSecurity,
		}

		kibanaRoleElasticsearchIndices[i] = kibanaRoleElasticsearchIndice
	}

	return kibanaRoleElasticsearchIndices
}

// buildRolesKibana permit to  build list of KibanaRoleKibana object
func buildRolesKibana(raws []interface{}) []kbapi.KibanaRoleKibana {
	kibanaRoleKibanas := make([]kbapi.KibanaRoleKibana, len(raws))

	for i, raw := range raws {
		m := raw.(map[string]interface{})

		kibanaRoleKibana := kbapi.KibanaRoleKibana{
			Base:    convertArrayInterfaceToArrayString(m["base"].(*schema.Set).List()),
			Feature: buildKibanaRoleKibanaFeatures(m["features"].(*schema.Set).List()),
			Spaces:  convertArrayInterfaceToArrayString(m["spaces"].(*schema.Set).List()),
		}

		kibanaRoleKibanas[i] = kibanaRoleKibana
	}

	return kibanaRoleKibanas
}

// buildKibanaRoleKibanaFeatures permit to build list of feature map
func buildKibanaRoleKibanaFeatures(raws []interface{}) map[string][]string {
	features := map[string][]string{}

	for _, raw := range raws {
		m := raw.(map[string]interface{})
		features[m["name"].(string)] = convertArrayInterfaceToArrayString(m["permissions"].(*schema.Set).List())
	}

	return features
}
