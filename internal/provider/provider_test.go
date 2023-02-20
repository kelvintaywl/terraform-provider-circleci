package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the HashiCups client is properly configured.
	// It is also possible to use the HASHICUPS_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "circleci" {
  // api_token via CIRCLE_TOKEN env var
  hostname = "circleci.com"
}
`
	// project name: github/kelvintaywl-cci/tf-provider-acceptance-test-dummy
	projectId   string = "c124cca6-d03e-4733-b84d-32b02347b78c"
	projectSlug string = "github/kelvintaywl-cci/tf-provider-acceptance-test-dummy"
	// webhook name: added-via-ui
	webhookId string = "8ed03fd1-5426-4138-a27d-aec0328c39fb"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"circleci": providerserver.NewProtocol6WithError(New()),
	}
)
