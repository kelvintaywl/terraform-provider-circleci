package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContextResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_context" "foobar_prod" {
	name         = "foobar_prod"
	owner        = {
		id = "%s"
		type = "organization"
	}
}
`, orgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_context.foobar_prod", "name", "foobar_prod"),
					resource.TestCheckResourceAttrSet("circleci_context.foobar_prod", "id"),
					resource.TestCheckResourceAttrSet("circleci_context.foobar_prod", "created_at"),
				),
			},
			// Test Import
			{
				ResourceName:        "circleci_context.foobar_prod",
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: fmt.Sprintf("organization,%s,", orgId),
			},
		},
	})
}
