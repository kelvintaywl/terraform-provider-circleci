resource "circleci_schedule" "every_first_day_of_quarter" {
  project_slug = "github/acmeorg/foobar"
  name         = "Quarterly schedule"
  description  = "Runs every 1st day of the quarter at 15:00~ UTC"
  branch       = "release"
  parameters = jsonencode({
    my_str  = "fizzbuzz"
    my_int  = 123
    my_bool = true
  })
  timetable = {
    per_hour      = 1
    hours_of_day  = [15]
    days_of_month = [1]
    months = [
      "JAN",
      "APR",
      "JUL",
      "OCT",
    ]
  }
  // refers to the owner of the CircleCI API token
  actor = "current"
}
