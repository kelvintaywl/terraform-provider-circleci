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
	// project name: github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
	projectId   string = "c124cca6-d03e-4733-b84d-32b02347b78c"
	projectSlug string = "github/kelvintaywl-tf/tf-provider-acceptance-test-dummy"
	// webhook name: added-via-ui
	webhookId string = "8ed03fd1-5426-4138-a27d-aec0328c39fb"

	// org slug: github/kelvintaywl-tf
	orgId string = "1e846a63-ae07-4549-a548-3db2aa4155e8"

	// context "from_tf" under github/kelvintaywl-tf org
	contextId   string = "c050c79f-f03b-4060-80b9-2562c7fdaa5c"
	contextName string = "from_tf"

	namespace     string = "kelvintaywl-tf"
	resourceClass string = "test"

	// standalone project & org data
	standaloneProjectId   string = "09cbbbea-993d-41fa-a467-57e1c543ead4"
	standaloneProjectSlug string = "circleci/7UQdtYSr1caLbAR2cHJdU7/2DACeEvUr7MosidActmnUs"

	// org circleci/7UQdtYSr1caLbAR2cHJdU7
	standaloneOrgId string = "346a7ade-9fae-47ec-b729-da3d5afbe4fc"

	// context "data-source-for-tf" under circleci/7UQdtYSr1caLbAR2cHJdU7
	standaloneContextId   string = "bb83cbf7-d20d-4224-9f70-3516111120a7"
	standaloneContextName string = "data-source-for-tf"
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
