resource "circleci_schedule" "weekly_schedule" {
  project_slug = "github/acmeorg/foobar"
  name         = "Weekly schedule"
  description  = "Runs every Monday at 00:00~ UTC"
  branch       = "main"
  timetable = {
    per_hour     = 1
    hours_of_day = [0]
    days_of_week = ["MON"]
    // we do not declare months here;
    // defaults to all months
  }
  // using the Scheduling system user
  actor = "system"
}
