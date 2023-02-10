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

data "circleci_webhooks" "project_webhooks" {
  // github/kelvintaywl-cci/ssh-ec2
  project_id = "a2502849-6bf0-486d-a357-75d331c65237"
}

output "webhooks" {
  description = "webhook_details"
  value       = data.circleci_webhooks.project_webhooks.webhooks == null ? null : data.circleci_webhooks.project_webhooks.webhooks[*]
}
