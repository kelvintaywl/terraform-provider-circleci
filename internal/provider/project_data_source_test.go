package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_project" "dummy" {
  slug = "%s"
}`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.circleci_project.dummy", "id", projectId),
					resource.TestCheckResourceAttr("data.circleci_project.dummy", "slug", projectSlug),

					resource.TestCheckResourceAttr("data.circleci_project.dummy", "organization_name", "kelvintaywl-tf"),
					resource.TestCheckResourceAttr("data.circleci_project.dummy", "organization_slug", "gh/kelvintaywl-tf"),
					resource.TestCheckResourceAttrSet("data.circleci_project.dummy", "organization_id"),

					resource.TestCheckResourceAttr("data.circleci_project.dummy", "vcs_url", "https://github.com/kelvintaywl-tf/tf-provider-acceptance-test-dummy"),
					resource.TestCheckResourceAttr("data.circleci_project.dummy", "vcs_default_branch", "main"),
					resource.TestCheckResourceAttr("data.circleci_project.dummy", "vcs_provider", "GitHub"),
				),
			},
		},
	})
}
