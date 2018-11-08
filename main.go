package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/spaceapegames/terraform-provider-example/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
