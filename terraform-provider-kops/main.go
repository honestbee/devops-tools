package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/honestbee/devops-tools/terraform-provider-kops/kops"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kops.Provider})
}
