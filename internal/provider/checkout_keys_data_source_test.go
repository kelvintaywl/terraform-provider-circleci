package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCheckoutKeysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_checkout_keys" "keys" {
  project_slug = "%s"
}`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.circleci_checkout_keys.keys", "id", projectSlug),
					resource.TestCheckResourceAttr("data.circleci_checkout_keys.keys", "project_slug", projectSlug),
				),
			},
		},
	})
}
