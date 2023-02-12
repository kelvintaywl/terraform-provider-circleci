package provider

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Prism returns example values from the underlying OpenAPI spec.
// See https://github.com/kelvintaywl/circleci-webhook-go-sdk/blob/main/openapi.yaml
const (
	// project name: github/kelvintaywl-cci/tf-provider-acceptance-test-dummy
	project_id string = "c124cca6-d03e-4733-b84d-32b02347b78c"
	// webhook name: added-via-ui
	webhook_id string = "4e6b7957-4448-489c-bdd5-6e6e0c37f15f"
)

func TestAccWebhooksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-cci/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_webhooks" "test" {
  project_id = "%s"
}`, project_id),
				Check: resource.ComposeAggregateTestCheckFunc(
					// top level
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "id", project_id),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.#", "1"),
					// webhook
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.id", webhook_id),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.name", "added-via-ui"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.verify_tls", "false"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.signing_secret", "****"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.id", project_id),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.type", "project"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.events.0", "workflow-completed"),
					resource.TestCheckResourceAttrSet("data.circleci_webhooks.test", "webhooks.0.created_at"),
					resource.TestCheckResourceAttrSet("data.circleci_webhooks.test", "webhooks.0.updated_at"),
				),
			},
		},
	})
}
