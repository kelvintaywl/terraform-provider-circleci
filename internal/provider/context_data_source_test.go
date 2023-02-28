package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContextDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-cci/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_context" "test" {
  name = "%s"
  owner = {
	id   = "%s"
	type = "organization"
  }
}`, "from_tf", orgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.circleci_context.test", "id"),
					resource.TestCheckResourceAttrSet("data.circleci_context.test", "created_at"),
				),
			},
		},
	})
}
