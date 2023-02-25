data "circleci_project" "my_project" {
  slug = "github/acmeorg/foobar"
}

resource "circleci_webhook" "my_webhook" {
  project_id     = data.circleci_project.my_project.id
  name           = "my_webhook"
  url            = "https://example.com/hook"
  signing_secret = "5uperSeCr3t!"
  verify_tls     = true
  events = [
    // accepts only "workflow-completed" and "job-completed"
    "job-completed",
    "workflow-completed"
  ]
}

output "webhooks" {
  description = "my_webhook_id"
  value       = circleci_webhook.my_webhook.id
}
