package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccProjectCheckoutKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_checkout_key" "my_key" {
	project_slug = "%s"
	type         = "deploy-key"
}
`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_checkout_key.my_key", "type", "deploy-key"),
					resource.TestCheckResourceAttr("circleci_checkout_key.my_key", "project_slug", projectSlug),
					resource.TestCheckResourceAttr("circleci_checkout_key.my_key", "preferred", "true"),
					resource.TestCheckResourceAttrSet("circleci_checkout_key.my_key", "fingerprint"),
					resource.TestCheckResourceAttrSet("circleci_checkout_key.my_key", "public_key"),
					resource.TestCheckResourceAttrSet("circleci_checkout_key.my_key", "id"),
				),
			},
			// Update and Read testing skipped; Checkout keys cannot be updated.
		},
	})
}
