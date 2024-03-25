package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccScheduleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_schedule" "my_schedule" {
	project_slug     = "%s"
	name             = "added-via-terraform-1"
	description      = "Runs weekly at 00:00~ every 1st of June, Dec"
	actor            = "current"
	branch           = "main"
	timetable = {
		per_hour      = 1
		hours_of_day  = [0]
		days_of_month = [1]
		months        = ["JUN", "DEC"]
	}
	parameters = jsonencode({
		my_int    = 123
		my_bool   = false
		my_string = "foobar"
	})
  }
`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "name", "added-via-terraform-1"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "project_slug", projectSlug),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "description", "Runs weekly at 00:00~ every 1st of June, Dec"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "actor", "current"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "branch", "main"),
					resource.TestCheckNoResourceAttr("circleci_schedule.my_schedule", "tag"),
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "parameters"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "timetable.per_hour", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.hours_of_day.*", "0"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.days_of_month.*", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.months.*", "JUN"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.months.*", "DEC"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "id"),
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "updated_at"),
				),
				ExpectNonEmptyPlan: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + fmt.Sprintf(`
			resource "circleci_schedule" "my_schedule" {
				project_slug     = "%s"
				name             = "added-via-terraform-2"
				description      = "Runs weekly at 00:00~ every 1st of June, Dec"
				actor            = "system"
				tag              = "v1.0"
				timetable = {
					per_hour      = 1
					hours_of_day  = [0]
					days_of_month = [1]
					months        = ["JUN", "DEC"]
				}
			  }
			`, projectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "name", "added-via-terraform-2"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "project_slug", projectSlug),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "description", "Runs weekly at 00:00~ every 1st of June, Dec"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "actor", "system"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "tag", "v1.0"),
					resource.TestCheckNoResourceAttr("circleci_schedule.my_schedule", "branch"),
					resource.TestCheckNoResourceAttr("circleci_schedule.my_schedule", "parameters"),
					resource.TestCheckResourceAttr("circleci_schedule.my_schedule", "timetable.per_hour", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.hours_of_day.*", "0"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.days_of_month.*", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.months.*", "JUN"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.my_schedule", "timetable.months.*", "DEC"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "id"),
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_schedule.my_schedule", "updated_at"),
				),
				ExpectNonEmptyPlan: true,
			},
			// Create and Read testing for standalone
			{
				Config: providerConfig + fmt.Sprintf(`
resource "circleci_schedule" "standalone" {
	project_slug     = "%s"
	name             = "added-via-terraform-1"
	description      = "Runs 2/hr at 00:00~03:00 every 1st & 8th of July"
	actor            = "current"
	branch           = "main"
	timetable = {
		per_hour      = 2
		hours_of_day  = [0,1,2]
		days_of_month = [1,8]
		months        = ["JUL"]
	}
	parameters = jsonencode({
		my_int    = 123
		my_bool   = false
		my_string = "foobar"
	})
  }
`, standaloneProjectSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("circleci_schedule.standalone", "name", "added-via-terraform-1"),
					resource.TestCheckResourceAttr("circleci_schedule.standalone", "project_slug", standaloneProjectSlug),
					resource.TestCheckResourceAttr("circleci_schedule.standalone", "description", "Runs 2/hr at 00:00~03:00 every 1st & 8th of July"),
					resource.TestCheckResourceAttr("circleci_schedule.standalone", "actor", "current"),
					resource.TestCheckResourceAttr("circleci_schedule.standalone", "branch", "main"),
					resource.TestCheckNoResourceAttr("circleci_schedule.standalone", "tag"),
					resource.TestCheckResourceAttrSet("circleci_schedule.standalone", "parameters"),
					resource.TestCheckResourceAttr("circleci_schedule.standalone", "timetable.per_hour", "2"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.standalone", "timetable.hours_of_day.*", "0"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.standalone", "timetable.hours_of_day.*", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.standalone", "timetable.hours_of_day.*", "2"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.standalone", "timetable.days_of_month.*", "1"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.standalone", "timetable.days_of_month.*", "8"),
					resource.TestCheckTypeSetElemAttr("circleci_schedule.standalone", "timetable.months.*", "JUL"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("circleci_schedule.standalone", "id"),
					resource.TestCheckResourceAttrSet("circleci_schedule.standalone", "created_at"),
					resource.TestCheckResourceAttrSet("circleci_schedule.standalone", "updated_at"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
