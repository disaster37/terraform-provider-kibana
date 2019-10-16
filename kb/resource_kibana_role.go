// Manage the role in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/role-management-api.html
// Supported version:
//  - v7
package kb

import (
	"fmt"

	kibana7 "github.com/disaster37/go-kibana-rest"
	kbapi7 "github.com/disaster37/go-kibana-rest/kbapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle user space in Kibana
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

// Create new user space in Kibana
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

	// Use the right client depend to Kibana version
	switch meta.(type) {
	// v7
	case *kibana7.Client:
		client := meta.(*kibana7.Client)

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
		d.Set("elasticsearch", []kbapi7.KibanaRoleElasticsearch{*role.Elasticsearch})
		d.Set("kibana", role.Kibana)
		d.Set("metadata", role.Metadata)
	default:
		return errors.New("Role is only supported by the kibana library >= v6!")
	}

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

	// Use the right client depend to Elasticsearch version
	switch meta.(type) {
	// v7
	case *kibana7.Client:
		client := meta.(*kibana7.Client)

		err := client.API.KibanaRoleManagement.Delete(id)
		if err != nil {
			if err.(kbapi7.APIError).Code == 404 {
				fmt.Printf("[WARN] Role %s not found - removing from state", id)
				log.Warnf("Role %s not found - removing from state", id)
				d.SetId("")
				return nil
			} else {
				return err
			}
		}

	default:
		return errors.New("Role is only supported by the kibana library >= v6!")
	}

	d.SetId("")

	log.Infof("Deleted role %s successfully", id)
	return nil

}

func createRole(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	metadataTemp := optionalInterfaceJson(d.Get("metadata").(string))
	elasticsearch := buildRolesElasticsearch(d.Get("elasticsearch").(*schema.Set).List())
	kibana := buildRolesKibana(d.Get("kibana").(*schema.Set).List())

	// Use the right client depend to Kibana version
	switch meta.(type) {
	// v7
	case *kibana7.Client:
		client := meta.(*kibana7.Client)

		var metadata map[string]interface{}
		if metadataTemp != nil {
			metadata = metadataTemp.(map[string]interface{})
		} else {
			metadata = nil
		}
		role := &kbapi7.KibanaRole{
			Name:          name,
			Elasticsearch: elasticsearch,
			Kibana:        kibana,
			Metadata:      metadata,
		}

		role, err := client.API.KibanaRoleManagement.CreateOrUpdate(role)
		if err != nil {
			return err
		}
	default:
		return errors.New("Role is only supported by the kibana library >= v6!")
	}

	return nil
}

func buildRolesElasticsearch(raws []interface{}) *kbapi7.KibanaRoleElasticsearch {
	if len(raws) == 0 {
		return nil
	}

	// We check only the first, we case use multiple KibanaRoleElasticsearch
	raw := raws[0].(map[string]interface{})

	kibanaRoleElasticsearch := &kbapi7.KibanaRoleElasticsearch{}

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

func buildKibanaRoleElasticsearchIndice(raws []interface{}) []kbapi7.KibanaRoleElasticsearchIndice {
	kibanaRoleElasticsearchIndices := make([]kbapi7.KibanaRoleElasticsearchIndice, len(raws))
	for i, raw := range raws {
		m := raw.(map[string]interface{})
		fieldSecurityTemp := optionalInterfaceJson(m["field_security"].(string))
		var fieldSecurity map[string]interface{}
		if fieldSecurityTemp != nil {
			fieldSecurity = fieldSecurityTemp.(map[string]interface{})
		} else {
			fieldSecurity = nil
		}
		kibanaRoleElasticsearchIndice := kbapi7.KibanaRoleElasticsearchIndice{
			Names:         convertArrayInterfaceToArrayString(m["names"].(*schema.Set).List()),
			Privileges:    convertArrayInterfaceToArrayString(m["privileges"].(*schema.Set).List()),
			Query:         optionalInterfaceJson(m["query"].(string)),
			FieldSecurity: fieldSecurity,
		}

		kibanaRoleElasticsearchIndices[i] = kibanaRoleElasticsearchIndice
	}

	return kibanaRoleElasticsearchIndices
}

func buildRolesKibana(raws []interface{}) []kbapi7.KibanaRoleKibana {
	kibanaRoleKibanas := make([]kbapi7.KibanaRoleKibana, len(raws))

	for i, raw := range raws {
		m := raw.(map[string]interface{})

		kibanaRoleKibana := kbapi7.KibanaRoleKibana{
			Base:    convertArrayInterfaceToArrayString(m["base"].(*schema.Set).List()),
			Feature: buildKibanaRoleKibanaFeatures(m["features"].(*schema.Set).List()),
			Spaces:  convertArrayInterfaceToArrayString(m["spaces"].(*schema.Set).List()),
		}

		kibanaRoleKibanas[i] = kibanaRoleKibana
	}

	return kibanaRoleKibanas
}

func buildKibanaRoleKibanaFeatures(raws []interface{}) map[string][]string {
	features := map[string][]string{}

	for _, raw := range raws {
		m := raw.(map[string]interface{})
		features[m["name"].(string)] = convertArrayInterfaceToArrayString(m["permissions"].(*schema.Set).List())
	}

	return features
}
