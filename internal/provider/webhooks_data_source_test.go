package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebhooksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_webhooks" "test" {
  project_id = "%s"
}`, projectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					// top level
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "id", projectId),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.#", "1"),
					// webhook
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.id", webhookId),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.name", "added-via-ui"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.url", "https://example.com/added-via-ui"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.verify_tls", "false"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.signing_secret", "****"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.id", projectId),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.type", "project"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.events.#", "2"),
					// TODO: make the below assertions work. somehow the parsing of the attr is not working here.
					// resource.TestCheckTypeSetElemAttr("data.circleci_webhooks.test", "webhooks.0.events.*", "webhook-completed"),
					// resource.TestCheckTypeSetElemAttr("data.circleci_webhooks.test", "webhooks.0.events.*", "job-completed"),
					resource.TestCheckResourceAttrSet("data.circleci_webhooks.test", "webhooks.0.created_at"),
					resource.TestCheckResourceAttrSet("data.circleci_webhooks.test", "webhooks.0.updated_at"),
				),
			},
		},
	})
}
