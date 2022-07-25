// Manage the role in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/role-management-api.html
// Supported version:
//  - v7

package kb

import (
	"encoding/json"
	"fmt"

	kibana "github.com/disaster37/go-kibana-rest/v7"
	kbapi "github.com/disaster37/go-kibana-rest/v7/kbapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
										DiffSuppressFunc: suppressEquivalentJSON,
									},
									"field_security": {
										Type:             schema.TypeString,
										Optional:         true,
										DiffSuppressFunc: suppressEquivalentJSON,
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
				MaxItems: 1,
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
				DiffSuppressFunc: suppressEquivalentJSON,
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
	fmt.Printf("[INFO] Created role %s successfully", name)

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
		log.Warnf("Role %s not found - removing from state", id)
		fmt.Printf("[WARN] Role %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	log.Debugf("Get role %s successfully:\n%s", id, role)

	d.Set("name", id)

	flattenKRE, err := flattenKibanaRoleElasticsearchMappings(role.Elasticsearch)
	if err != nil {
		return err
	}

	log.Debugf("Flatten ES: +%v\n", flattenKRE)
	if err := d.Set("elasticsearch", flattenKRE); err != nil {
		return fmt.Errorf("error setting elasticsearch: %w", err)
	}

	if err := d.Set("kibana", flattenKibanaRoleKibanaMappings(role.Kibana)); err != nil {
		return fmt.Errorf("error setting kibana: %w", err)
	}

	flattenKRM, err := convertInterfaceToJsonString(role.Metadata)
	if err != nil {
		return err
	}
	if err := d.Set("metadata", flattenKRM); err != nil {
		return fmt.Errorf("error setting metadata: %w", err)
	}

	log.Infof("Read role %s successfully", id)
	fmt.Printf("[INFO] Read role %s successfully", id)

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
	fmt.Printf("[INFO] Updated role %s successfully", id)

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
			log.Warnf("Role %s not found - removing from state", id)
			fmt.Printf("[WARN] Role %s not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return err

	}

	d.SetId("")

	log.Infof("Deleted role %s successfully", id)
	fmt.Printf("[INFO] Deleted role %s successfully", id)
	return nil

}

// createRole permit to create or update role in Kibana
func createRole(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	metadataTemp := optionalInterfaceJSON(d.Get("metadata").(string))
	roleElasticsearch, err := buildRolesElasticsearch(d.Get("elasticsearch").(*schema.Set).List())
	if err != nil {
		return err
	}
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

	role, err = client.API.KibanaRoleManagement.CreateOrUpdate(role)
	if err != nil {
		return err
	}

	return nil
}

// buildRolesElasticsearch permit to construct kibanaRoleElasticsearch object
func buildRolesElasticsearch(raws []interface{}) (*kbapi.KibanaRoleElasticsearch, error) {
	if len(raws) == 0 {
		return nil, nil
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
		krei, err := buildKibanaRoleElasticsearchIndice(raw["indices"].(*schema.Set).List())
		if err != nil {
			return nil, err
		}
		kibanaRoleElasticsearch.Indices = krei
	}

	return kibanaRoleElasticsearch, nil

}

// buildKibanaRoleElasticsearchIndice permit to build list of KibanaRoleElasticsearchIndice
func buildKibanaRoleElasticsearchIndice(raws []interface{}) ([]kbapi.KibanaRoleElasticsearchIndice, error) {
	kibanaRoleElasticsearchIndices := make([]kbapi.KibanaRoleElasticsearchIndice, len(raws))
	for i, raw := range raws {
		m := raw.(map[string]interface{})
		fieldSecurityTemp := optionalInterfaceJSON(m["field_security"].(string))
		var fieldSecurity map[string]interface{}
		if fieldSecurityTemp != nil {

			if err := json.Unmarshal(fieldSecurityTemp.(json.RawMessage), &fieldSecurity); err != nil {
				return nil, err
			}

		} else {
			fieldSecurity = nil
		}
		kibanaRoleElasticsearchIndice := kbapi.KibanaRoleElasticsearchIndice{
			Names:         convertArrayInterfaceToArrayString(m["names"].(*schema.Set).List()),
			Privileges:    convertArrayInterfaceToArrayString(m["privileges"].(*schema.Set).List()),
			Query:         m["query"].(string),
			FieldSecurity: fieldSecurity,
		}

		kibanaRoleElasticsearchIndices[i] = kibanaRoleElasticsearchIndice
	}

	return kibanaRoleElasticsearchIndices, nil
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

func flattenKibanaRoleElasticsearchMappings(kre *kbapi.KibanaRoleElasticsearch) ([]interface{}, error) {

	// Handle empty object
	if kre == nil || (len(kre.Cluster) == 0 && len(kre.Indices) == 0 && len(kre.RunAs) == 0) {
		return nil, nil
	}

	var tfList []interface{}
	flattenKRE, err := flattenKibanaRoleElasticsearchMapping(kre)
	if err != nil {
		return nil, err
	}
	tfList = append(tfList, flattenKRE)

	return tfList, nil
}

func flattenKibanaRoleElasticsearchMapping(kre *kbapi.KibanaRoleElasticsearch) (map[string]interface{}, error) {
	if kre == nil {
		return nil, nil
	}

	tfMap := make(map[string]interface{})

	if kre.Indices != nil {
		flatten, err := flattenKibanaRoleElasticsearchMappingsIndices(kre.Indices)
		if err != nil {
			return nil, err
		}
		tfMap["indices"] = flatten
	}

	if kre.Cluster != nil {
		tfMap["cluster"] = kre.Cluster
	} else {
		tfMap["cluster"] = make([]interface{}, 0)
	}

	if kre.RunAs != nil {
		tfMap["run_as"] = kre.RunAs
	} else {
		tfMap["run_as"] = make([]interface{}, 0)
	}

	return tfMap, nil
}

func flattenKibanaRoleElasticsearchMappingsIndices(krei []kbapi.KibanaRoleElasticsearchIndice) ([]interface{}, error) {
	if krei == nil {
		return nil, nil
	}

	tfList := make([]interface{}, 0)

	for _, item := range krei {
		flatten, err := flattenKibanaRoleElasticsearchMappingIndices(item)
		if err != nil {
			return nil, err
		}
		tfList = append(tfList, flatten)
	}

	return tfList, nil
}

func flattenKibanaRoleElasticsearchMappingIndices(krei kbapi.KibanaRoleElasticsearchIndice) (map[string]interface{}, error) {

	tfMap := make(map[string]interface{})

	tfMap["names"] = krei.Names
	tfMap["privileges"] = krei.Privileges
	flattenQuerry, err := convertInterfaceToJsonString(krei.Query)
	tfMap["query"] = flattenQuerry

	flattenFieldSecurity, err := convertInterfaceToJsonString(krei.FieldSecurity)
	if err != nil {
		return nil, err
	}
	tfMap["field_security"] = flattenFieldSecurity

	return tfMap, nil

}

func flattenKibanaRoleFeatureMappings(krf map[string][]string) []interface{} {
	if krf == nil {
		return nil
	}

	tfList := make([]interface{}, 0)

	for name, item := range krf {
		tfMap := make(map[string]interface{})
		tfMap["name"] = name
		tfMap["permissions"] = item
		tfList = append(tfList, tfMap)
	}

	return tfList
}

func flattenKibanaRoleKibanaMapping(krk kbapi.KibanaRoleKibana) map[string]interface{} {
	tfMap := make(map[string]interface{})
	tfMap["base"] = krk.Base
	tfMap["spaces"] = krk.Spaces
	tfMap["features"] = flattenKibanaRoleFeatureMappings(krk.Feature)

	return tfMap
}

func flattenKibanaRoleKibanaMappings(krk []kbapi.KibanaRoleKibana) []interface{} {
	if krk == nil {
		return nil
	}

	tfList := make([]interface{}, 0, len(krk))

	for _, item := range krk {
		tfList = append(tfList, flattenKibanaRoleKibanaMapping(item))
	}

	return tfList
}
