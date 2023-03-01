package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContextEnvVarResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
data "circleci_context" "test" {
	name = "%s"
	owner = {
		id   = "%s"
		type = "organization"
	}
}

resource "circleci_context_env_var" "env1" {
	name         = "FOOBAR"
	value        = "random1234"
	context_id   = data.circleci_context.test.id
}
`, contextName, orgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_context_env_var.env1", "name", "FOOBAR"),
					resource.TestCheckResourceAttr("circleci_context_env_var.env1", "value", "random1234"),
					resource.TestCheckResourceAttr("circleci_context_env_var.env1", "id", fmt.Sprintf("%s/FOOBAR", contextId)),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env1", "context_id"),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env1", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env1", "updated_at"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
data "circleci_context" "test" {
	name = "%s"
	owner = {
		id   = "%s"
		type = "organization"
	}
}

resource "circleci_context_env_var" "env1" {
	name         = "FOOBAR"
	value        = "changed"
	context_id   = data.circleci_context.test.id
}

resource "circleci_context_env_var" "env2" {
	name         = "FIZZBUZZ"
	value        = "Lorem Ipsum"
	context_id   = data.circleci_context.test.id
}
`, contextName, orgId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_context_env_var.env1", "name", "FOOBAR"),
					resource.TestCheckResourceAttr("circleci_context_env_var.env1", "value", "changed"),
					resource.TestCheckResourceAttr("circleci_context_env_var.env1", "id", fmt.Sprintf("%s/FOOBAR", contextId)),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env1", "context_id"),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env1", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env1", "updated_at"),

					resource.TestCheckResourceAttr("circleci_context_env_var.env2", "name", "FIZZBUZZ"),
					resource.TestCheckResourceAttr("circleci_context_env_var.env2", "value", "Lorem Ipsum"),
					resource.TestCheckResourceAttr("circleci_context_env_var.env2", "id", fmt.Sprintf("%s/FIZZBUZZ", contextId)),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env2", "context_id"),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env2", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_context_env_var.env2", "updated_at"),
				),
			},
		},
	})
}
