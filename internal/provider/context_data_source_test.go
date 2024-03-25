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
				// github/kelvintaywl-tf
				Config: providerConfig + fmt.Sprintf(`
data "circleci_context" "test" {
  name = "%s"
  owner = {
	id   = "%s"
	type = "organization"
  }
}`, contextName, orgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.circleci_context.test", "id", contextId),
					resource.TestCheckResourceAttrSet("data.circleci_context.test", "created_at"),
				),
			},
			// Read testing for standalone
			{
				// circleci/7UQdtYSr1caLbAR2cHJdU7
				Config: providerConfig + fmt.Sprintf(`
data "circleci_context" "standalone" {
  name = "%s"
  owner = {
	id   = "%s"
	type = "organization"
  }
}`, standaloneContextName, standaloneOrgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.circleci_context.standalone", "id", standaloneContextId),
					resource.TestCheckResourceAttrSet("data.circleci_context.standalone", "created_at"),
				),
			},
		},
	})
}
