// Return the connection settings of Kibana
// Supported version:
//  - v7

package kb

import (
	kibana "github.com/disaster37/go-kibana-rest/v7"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceKibanaHost() *schema.Resource {
	return &schema.Resource{
		Description: "`kibana_host` can be used to retrieve the Kibana connection settings.",
		Read:        dataSourceKibanaHostRead,

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

func dataSourceKibanaHostRead(d *schema.ResourceData, m interface{}) error {
	var url string
	var username string
	var password string

	conf := m.(*kibana.Client)

	url = conf.Client.HostURL
	username = conf.Client.UserInfo.Username
	password = conf.Client.UserInfo.Password

	d.SetId(url)
	d.Set("url", url)
	d.Set("username", username)
	d.Set("password", password)

	return nil
}
