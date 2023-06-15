package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRunnerTokenResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_runner_token" "delete_me" {
	resource_class = "%s/test"
	nickname       = "acceptance-test"
}
`, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_runner_token.delete_me", "resource_class", "kelvintaywl-tf/test"),
					resource.TestCheckResourceAttr("circleci_runner_token.delete_me", "nickname", "acceptance-test"),
					resource.TestCheckResourceAttrSet("circleci_runner_token.delete_me", "id"),
					resource.TestCheckResourceAttrSet("circleci_runner_token.delete_me", "token"),
				),
			},
		},
	})
}
