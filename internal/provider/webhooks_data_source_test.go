package provider

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Prism returns example values from the underlying OpenAPI spec.
// See https://github.com/kelvintaywl/circleci-webhook-go-sdk/blob/main/openapi.yaml
const (
	project_id         string = "c124cca6-d03e-4733-b84d-32b02347b78c"
	webhook_id         string = "d57ecc67-7a3b-4fd9-a1b4-442d4703bb8d"
	created_updated_at string = "2023-02-10T04:49:36.117Z"
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
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.name", "test1"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.verify_tls", "true"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.signing_secret", "****"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.id", project_id),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.scope.type", "project"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.events.0", "workflow-completed"),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.created_at", created_updated_at),
					resource.TestCheckResourceAttr("data.circleci_webhooks.test", "webhooks.0.updated_at", created_updated_at),
				),
			},
		},
	})
}
