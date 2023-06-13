package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRunnerResourceClassesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				// github/kelvintaywl-tf/tf-provider-acceptance-test-dummy
				Config: providerConfig + fmt.Sprintf(`
data "circleci_runner_resource_classes" "test" {
  namespace = "%s"
}`, namespace),
				Check: resource.ComposeAggregateTestCheckFunc(
					// top level
					resource.TestCheckResourceAttr("data.circleci_runner_resource_classes.test", "id", namespace),
					resource.TestCheckResourceAttr("data.circleci_runner_resource_classes.test", "resource_classes.#", "1"),
					// resource-class
					resource.TestCheckResourceAttr("data.circleci_runner_resource_classes.test", "resource_classes.0.id", "1a24212b-7db4-493d-9469-1785e99af123"),
					resource.TestCheckResourceAttr("data.circleci_runner_resource_classes.test", "resource_classes.0.resource_class", "kelvintaywl-tf/test"),
					resource.TestCheckResourceAttr("data.circleci_runner_resource_classes.test", "resource_classes.0.description", "throwaway runner"),
				),
			},
		},
	})
}
