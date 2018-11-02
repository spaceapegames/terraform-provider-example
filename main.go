package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/spaceapegames/terraform-provider-blog/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc:    provider.Provider,
	})
}
