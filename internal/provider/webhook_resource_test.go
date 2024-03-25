package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWebhookResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_webhook" "my_webhook" {
	project_id     = "%s"
	name           = "added-via-terraform-1"
	url            = "https://example.com/added-via-terraform"
	signing_secret = "rand0m5eCr3t"
	verify_tls     = true
	events = [
	  "job-completed"
	]
}
`, projectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "name", "added-via-terraform-1"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "url", "https://example.com/added-via-terraform"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "signing_secret", "rand0m5eCr3t"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "project_id", projectId),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "verify_tls", "true"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "events.#", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_webhook.my_webhook", "events.*", "job-completed"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("circleci_webhook.my_webhook", "id"),
					resource.TestCheckResourceAttrSet("circleci_webhook.my_webhook", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_webhook.my_webhook", "updated_at"),
				),
				// ExpectNonEmptyPlan: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_webhook" "my_webhook" {
	project_id     = "%s"
	name           = "added-via-terraform-1"
	url            = "https://example.com/added-via-terraform"
	signing_secret = "changed"
	verify_tls     = false
	events = [
	  "workflow-completed"
	]
}
`, projectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "name", "added-via-terraform-1"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "url", "https://example.com/added-via-terraform"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "signing_secret", "changed"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "project_id", projectId),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "verify_tls", "false"),
					resource.TestCheckResourceAttr("circleci_webhook.my_webhook", "events.0", "workflow-completed"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("circleci_webhook.my_webhook", "id"),
					resource.TestCheckResourceAttrSet("circleci_webhook.my_webhook", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_webhook.my_webhook", "updated_at"),
				),
				// ExpectNonEmptyPlan: true,
			},
			// No updates
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_webhook" "my_webhook" {
    project_id     = "%s"
    name           = "added-via-terraform-1"
    url            = "https://example.com/added-via-terraform"
    signing_secret = "changed"
    verify_tls     = false
    events = [
      "workflow-completed"
    ]
}
`, projectId),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Test Import
			{
				ResourceName: "circleci_webhook.my_webhook",
				ImportState:  true,
				// signing_secret will return masked
				ImportStateVerifyIgnore: []string{"signing_secret"},
			},
			// Create and Read testing for standalone
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_webhook" "standalone" {
	project_id     = "%s"
	name           = "added-via-terraform-1"
	url            = "https://standalone.example.com/added-via-terraform"
	signing_secret = "st4nD4L0n3"
	verify_tls     = true
	events = [
	  "job-completed"
	]
}
`, standaloneProjectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_webhook.standalone", "name", "added-via-terraform-1"),
					resource.TestCheckResourceAttr("circleci_webhook.standalone", "url", "https://standalone.example.com/added-via-terraform"),
					resource.TestCheckResourceAttr("circleci_webhook.standalone", "signing_secret", "st4nD4L0n3"),
					resource.TestCheckResourceAttr("circleci_webhook.standalone", "project_id", standaloneProjectId),
					resource.TestCheckResourceAttr("circleci_webhook.standalone", "verify_tls", "true"),
					resource.TestCheckResourceAttr("circleci_webhook.standalone", "events.#", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_webhook.standalone", "events.*", "job-completed"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("circleci_webhook.standalone", "id"),
					resource.TestCheckResourceAttrSet("circleci_webhook.standalone", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_webhook.standalone", "updated_at"),
				),
				// ExpectNonEmptyPlan: true,
			},
		},
	})
}
