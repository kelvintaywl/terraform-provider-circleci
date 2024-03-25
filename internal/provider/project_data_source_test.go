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
			// Read testing for standalone
			{
				// github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_project" "ipsum" {
  slug = "%s"
}`, standaloneProjectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "id", standaloneProjectId),
					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "slug", standaloneProjectSlug),
					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "name", "loren-ipsum"),

					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "organization_name", "ktwl41"),
					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "organization_slug", "circleci/7UQdtYSr1caLbAR2cHJdU7"),
					resource.TestCheckResourceAttrSet("data.circleci_project.ipsum", "organization_id"),

					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "vcs_url", "//circleci.com/346a7ade-9fae-47ec-b729-da3d5afbe4fc/09cbbbea-993d-41fa-a467-57e1c543ead4"),
					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "vcs_default_branch", "main"),
					resource.TestCheckResourceAttr("data.circleci_project.ipsum", "vcs_provider", "CircleCI"),
				),
			},
		},
	})
}
