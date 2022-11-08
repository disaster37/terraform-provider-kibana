// Return the connection settings of Kibana
// Supported version:
//  - v7

package kb

import (
	"context"

	kibana "github.com/disaster37/go-kibana-rest/v8"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKibanaHost() *schema.Resource {
	return &schema.Resource{
		Description: "`kibana_host` can be used to retrieve the Kibana connection settings.",
		ReadContext: dataSourceKibanaHostRead,

		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Kibana URL",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Username to use to connect to Kibana using basic auth",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Password to use to connect to Kibana using basic auth",
			},
		},
	}
}

func dataSourceKibanaHostRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var url string
	var username string
	var password string
	var err error

	conf := m.(*kibana.Client)

	url = conf.Client.HostURL
	username = conf.Client.UserInfo.Username
	password = conf.Client.UserInfo.Password

	d.SetId(url)
	if err = d.Set("url", url); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("username", username); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("password", password); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
