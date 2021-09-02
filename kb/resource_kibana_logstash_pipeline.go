// Manage the logstash pipeline in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/logstash-configuration-management-api.html
// Supported version:
//  - v7

package kb

import (
	kibana "github.com/disaster37/go-kibana-rest/v7"
	kbapi "github.com/disaster37/go-kibana-rest/v7/kbapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle logstash pipeline in Kibana
func resourceKibanaLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceKibanaLogstashPipelineCreate,
		Read:   resourceKibanaLogstashPipelineRead,
		Update: resourceKibanaLogstashPipelineUpdate,
		Delete: resourceKibanaLogstashPipelineDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipeline": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressEquivalentJSON,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// Create new logstash pipeline in Kibana
func resourceKibanaLogstashPipelineCreate(d *schema.ResourceData, meta interface{}) error {

	logstashPipeline, err := createOrUpdateLogstashPipeline(d, meta)
	if err != nil {
		return err
	}

	d.SetId(logstashPipeline.ID)

	log.Infof("Created logstash pipeline %s successfully", logstashPipeline.ID)

	return resourceKibanaLogstashPipelineRead(d, meta)
}

// Read existing logstash pipeline in Kibana
func resourceKibanaLogstashPipelineRead(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()

	log.Debugf("Logstash pipeline id:  %s", id)

	client := meta.(*kibana.Client)

	logstashPiepeline, err := client.API.KibanaLogstashPipeline.Get(id)
	if err != nil {
		return err
	}

	if logstashPiepeline == nil {
		log.Warnf("Logstash piepline %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	log.Debugf("Get logstash piepeline %s successfully:\n%s", id, logstashPiepeline)

	d.Set("name", logstashPiepeline.ID)
	d.Set("description", logstashPiepeline.Description)
	d.Set("username", logstashPiepeline.Username)
	d.Set("pipeline", logstashPiepeline.Pipeline)
	d.Set("settings", logstashPiepeline.Settings)

	log.Infof("Read logstash pipeline %s successfully", id)

	return nil
}

// Update existing logstash pipeline in Elasticsearch
func resourceKibanaLogstashPipelineUpdate(d *schema.ResourceData, meta interface{}) error {

	logstashPipeline, err := createOrUpdateLogstashPipeline(d, meta)
	if err != nil {
		return err
	}

	log.Infof("Updated logstash piepeline %s successfully", logstashPipeline.ID)

	return resourceKibanaLogstashPipelineRead(d, meta)
}

// Delete existing logstash pipeline in Elasticsearch
func resourceKibanaLogstashPipelineDelete(d *schema.ResourceData, meta interface{}) error {

	id := d.Id()
	log.Debugf("Logstash pipeline id: %s", id)

	client := meta.(*kibana.Client)

	err := client.API.KibanaLogstashPipeline.Delete(id)
	if err != nil {
		if err.(kbapi.APIError).Code == 404 {
			log.Warnf("Logstash pipeline %s not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return err

	}

	d.SetId("")

	log.Infof("Deleted logstash pipeline %s successfully", id)
	return nil

}

// createOrUpdateLogstashPipeline permit to create or update logstash pipeline
func createOrUpdateLogstashPipeline(d *schema.ResourceData, meta interface{}) (*kbapi.LogstashPipeline, error) {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	pipeline := d.Get("pipeline").(string)
	settings := d.Get("settings").(map[string]interface{})

	client := meta.(*kibana.Client)

	logstashPipeline := &kbapi.LogstashPipeline{
		ID:          name,
		Description: description,
		Pipeline:    pipeline,
		Settings:    settings,
	}

	logstashPipeline, err := client.API.KibanaLogstashPipeline.CreateOrUpdate(logstashPipeline)
	if err != nil {
		return nil, err
	}

	return logstashPipeline, nil
}
