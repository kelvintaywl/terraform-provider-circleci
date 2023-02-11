terraform {
  required_providers {
    circleci = {
      source = "example.com/kelvintaywl/circleci"
    }
  }
}

provider "circleci" {
  // api_token = ""; or set via CIRCLE_TOKEN
  // hostname = "https://circleci.com"; or set via CIRCLE_HOSTNAME
}

locals {
  // Replace this with your CircleCI project ID
  project_id = "c124cca6-d03e-4733-b84d-32b02347b78c"
}

resource "circleci_webhook" "my_webhook" {
  project_id     = local.project_id
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
