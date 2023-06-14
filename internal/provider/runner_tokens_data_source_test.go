package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRunnerTokensDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_runner_tokens" "test" {
  resource_class = "%s/%s"
}`, namespace, resourceClass),
				Check: resource.ComposeAggregateTestCheckFunc(
					// top level
					resource.TestCheckResourceAttr("data.circleci_runner_tokens.test", "id", "kelvintaywl-tf/test"),
					resource.TestCheckResourceAttr("data.circleci_runner_tokens.test", "tokens.#", "1"),
					// token
					resource.TestCheckResourceAttr("data.circleci_runner_tokens.test", "tokens.0.id", "0912b935-3026-45df-97e4-dda04154a7d2"),
					resource.TestCheckResourceAttr("data.circleci_runner_tokens.test", "tokens.0.resource_class", "kelvintaywl-tf/test"),
					resource.TestCheckResourceAttr("data.circleci_runner_tokens.test", "tokens.0.nickname", "default"),
					resource.TestCheckResourceAttrSet("data.circleci_runner_tokens.test", "tokens.0.created_at"),
				),
			},
		},
	})
}
