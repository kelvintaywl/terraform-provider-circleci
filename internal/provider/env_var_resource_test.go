package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectEnvVarResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_env_var" "env1" {
	project_slug = "%s"
	name         = "FOOBAR"
	value        = "random1234"
}
`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_env_var.env1", "name", "FOOBAR"),
					resource.TestCheckResourceAttr("circleci_env_var.env1", "project_slug", projectSlug),
					resource.TestCheckResourceAttr("circleci_env_var.env1", "value", "random1234"),
					resource.TestCheckResourceAttr("circleci_env_var.env1", "id", fmt.Sprintf("%s/FOOBAR", projectSlug)),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_env_var" "env1" {
	project_slug = "%s"
	name         = "FOOBAR"
	value        = "random1234"
}

resource "circleci_env_var" "env2" {
	project_slug = "%s"
	name         = "FIZZBUZZ"
	value        = "Lorem Ipsum"
}
`, projectSlug, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_env_var.env1", "name", "FOOBAR"),
					resource.TestCheckResourceAttr("circleci_env_var.env1", "project_slug", projectSlug),
					resource.TestCheckResourceAttr("circleci_env_var.env1", "value", "random1234"),
					resource.TestCheckResourceAttr("circleci_env_var.env1", "id", fmt.Sprintf("%s/%s", projectSlug, "FOOBAR")),

					resource.TestCheckResourceAttr("circleci_env_var.env2", "name", "FIZZBUZZ"),
					resource.TestCheckResourceAttr("circleci_env_var.env2", "project_slug", projectSlug),
					resource.TestCheckResourceAttr("circleci_env_var.env2", "value", "Lorem Ipsum"),
					resource.TestCheckResourceAttr("circleci_env_var.env2", "id", fmt.Sprintf("%s/FIZZBUZZ", projectSlug)),
				),
			},
		},
	})
}
