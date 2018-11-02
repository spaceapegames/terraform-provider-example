package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spaceapegames/terraform-provider-blog/api/client"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_ADDRESS", ""),
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_PORT", ""),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"example_item": resourceItem(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	address := d.Get("SERVICE_ADDRESS").(string)
	port := d.Get("SERVICE_ADDRESS").(int)
	token := d.Get("SERVICE_TOKEN").(string)
	return client.NewClient(address, port, token), nil

}