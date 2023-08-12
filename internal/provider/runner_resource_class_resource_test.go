package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRunnerResourceClassResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_runner_resource_class" "blade_runner" {
	resource_class = "%s/acceptance-test"
	description    = "From Terraform acceptance test"
}
`, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_runner_resource_class.blade_runner", "resource_class", "kelvintaywl-tf/acceptance-test"),
					resource.TestCheckResourceAttr("circleci_runner_resource_class.blade_runner", "description", "From Terraform acceptance test"),
					resource.TestCheckResourceAttrSet("circleci_runner_resource_class.blade_runner", "id"),
				),
			},
			// Test Import
			{
				ResourceName:        "circleci_runner_resource_class.blade_runner",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: fmt.Sprintf("%s/acceptance-test,", namespace),
			},
		},
	})
}
