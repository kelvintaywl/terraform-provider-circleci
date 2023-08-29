package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_project" "p1" {
	slug = "%s"
}
`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_project.p1", "slug", projectSlug),

					resource.TestCheckResourceAttr("circleci_project.p1", "organization_name", "kelvintaywl-tf"),
					resource.TestCheckResourceAttr("circleci_project.p1", "organization_slug", "gh/kelvintaywl-tf"),
					resource.TestCheckResourceAttrSet("circleci_project.p1", "organization_id"),

					resource.TestCheckResourceAttr("circleci_project.p1", "vcs_url", "https://github.com/kelvintaywl-tf/tf-provider-acceptance-test-dummy"),
					resource.TestCheckResourceAttr("circleci_project.p1", "vcs_default_branch", "main"),
					resource.TestCheckResourceAttr("circleci_project.p1", "vcs_provider", "GitHub"),
				),
			},
		},
	})
}
