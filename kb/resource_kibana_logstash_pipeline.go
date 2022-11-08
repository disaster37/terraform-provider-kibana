// Manage the logstash pipeline in Kibana
// API documentation: https://www.elastic.co/guide/en/kibana/master/logstash-configuration-management-api.html
// Supported version:
//  - v7

package kb

import (
	"context"
	"fmt"

	kibana "github.com/disaster37/go-kibana-rest/v8"
	kbapi "github.com/disaster37/go-kibana-rest/v8/kbapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Resource specification to handle logstash pipeline in Kibana
func resourceKibanaLogstashPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKibanaLogstashPipelineCreate,
		ReadContext:   resourceKibanaLogstashPipelineRead,
		UpdateContext: resourceKibanaLogstashPipelineUpdate,
		DeleteContext: resourceKibanaLogstashPipelineDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pipeline_workers": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pipeline_batch_size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pipeline_batch_delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pipeline_ecs_compatibility": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pipeline_ordored": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue_max_bytes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue_checkpoint_writes": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

// Create new logstash pipeline in Kibana
func resourceKibanaLogstashPipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	logstashPipeline, err := createOrUpdateLogstashPipeline(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(logstashPipeline.ID)

	log.Infof("Created logstash pipeline %s successfully", logstashPipeline.ID)
	fmt.Printf("[INFO] Created logstash pipeline %s successfully", logstashPipeline.ID)

	return resourceKibanaLogstashPipelineRead(ctx, d, meta)
}

// Read existing logstash pipeline in Kibana
func resourceKibanaLogstashPipelineRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var err error
	id := d.Id()

	log.Debugf("Logstash pipeline id:  %s", id)

	client := meta.(*kibana.Client)

	logstashPiepeline, err := client.API.KibanaLogstashPipeline.Get(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if logstashPiepeline == nil {
		log.Warnf("Logstash piepline %s not found - removing from state", id)
		fmt.Printf("[WARN] Logstash piepline %s not found - removing from state", id)
		d.SetId("")
		return nil
	}

	log.Debugf("Get logstash piepeline %s successfully:\n%s", id, logstashPiepeline)

	if err = d.Set("name", logstashPiepeline.ID); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("description", logstashPiepeline.Description); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("username", logstashPiepeline.Username); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("pipeline", logstashPiepeline.Pipeline); err != nil {
		return diag.FromErr(err)
	}

	if len(logstashPiepeline.Settings) > 0 {

		settings := make([]any, 0, 1)
		setting := map[string]any{
			"pipeline_workers":           logstashPiepeline.Settings["pipeline.workers"],
			"pipeline_batch_size":        logstashPiepeline.Settings["pipeline.batch.size"],
			"pipeline_batch_delay":       logstashPiepeline.Settings["pipeline.batch.delay"],
			"pipeline_ecs_compatibility": logstashPiepeline.Settings["pipeline.ecs_compatibility"],
			"pipeline_ordored":           logstashPiepeline.Settings["pipeline.ordered"],
			"queue_type":                 logstashPiepeline.Settings["queue.type"],
			"queue_max_bytes":            logstashPiepeline.Settings["queue.max_bytes"],
			"queue_checkpoint_writes":    logstashPiepeline.Settings["queue.checkpoint.writes"],
		}

		settings = append(settings, setting)
		if err = d.Set("settings", settings); err != nil {
			return diag.FromErr(err)
		}

	}

	log.Infof("Read logstash pipeline %s successfully", id)
	fmt.Printf("[INFO] Read logstash pipeline %s successfully", id)

	return nil
}

// Update existing logstash pipeline in Elasticsearch
func resourceKibanaLogstashPipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	logstashPipeline, err := createOrUpdateLogstashPipeline(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Infof("Updated logstash piepeline %s successfully", logstashPipeline.ID)
	fmt.Printf("[INFO] Updated logstash piepeline %s successfully", logstashPipeline.ID)

	return resourceKibanaLogstashPipelineRead(ctx, d, meta)
}

// Delete existing logstash pipeline in Elasticsearch
func resourceKibanaLogstashPipelineDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	id := d.Id()
	log.Debugf("Logstash pipeline id: %s", id)

	client := meta.(*kibana.Client)

	if err := client.API.KibanaLogstashPipeline.Delete(id); err != nil {
		if err.(kbapi.APIError).Code == 404 {
			log.Warnf("Logstash pipeline %s not found - removing from state", id)
			fmt.Printf("[WARN] Logstash pipeline %s not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)

	}

	d.SetId("")

	log.Infof("Deleted logstash pipeline %s successfully", id)
	fmt.Printf("[INFO] Deleted logstash pipeline %s successfully", id)
	return nil

}

// createOrUpdateLogstashPipeline permit to create or update logstash pipeline
func createOrUpdateLogstashPipeline(d *schema.ResourceData, meta interface{}) (*kbapi.LogstashPipeline, error) {
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	pipeline := d.Get("pipeline").(string)
	settings := d.Get("settings").(*schema.Set).List()

	client := meta.(*kibana.Client)

	logstashPipeline := &kbapi.LogstashPipeline{
		ID:          name,
		Description: description,
		Pipeline:    pipeline,
		Settings:    map[string]any{},
	}

	if len(settings) > 0 {
		for key, value := range settings[0].(map[string]any) {
			switch key {
			case "pipeline_workers":
				logstashPipeline.Settings["pipeline.workers"] = value
			case "pipeline_batch_size":
				logstashPipeline.Settings["pipeline.batch.size"] = value
			case "pipeline_batch_delay":
				logstashPipeline.Settings["pipeline.batch.delay"] = value
			case "pipeline_ecs_compatibility":
				logstashPipeline.Settings["pipeline.ecs_compatibility"] = value
			case "pipeline_ordored":
				logstashPipeline.Settings["pipeline.ordered"] = value
			case "queue_type":
				logstashPipeline.Settings["queue.type"] = value
			case "queue_max_bytes":
				logstashPipeline.Settings["queue.max_bytes"] = value
			case "queue_checkpoint_writes":
				logstashPipeline.Settings["queue.checkpoint.writes"] = value
			}
		}
	}

	logstashPipeline, err := client.API.KibanaLogstashPipeline.CreateOrUpdate(logstashPipeline)
	if err != nil {
		return nil, err
	}

	return logstashPipeline, nil
}
