package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebhooksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-cci/tf-provider-acceptance-test-dummy
				Config: providerConfig + `
data "circleci_webhooks" "test" {
  project_id = "c124cca6-d03e-4733-b84d-32b02347b78c"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// top level
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "id", "c124cca6-d03e-4733-b84d-32b02347b78c"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.#", "1"),
					// webhook
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.id", "string"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.name", "string"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.verify_tls", "true"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.signing_secret", "string"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.id", "string"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.type", "project"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.events.0", "workflow-completed"),
				),
			},
		},
	})
}
