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
  project_id = "32bbe47f-2bdf-4bb7-8390-ce682161a95f"
}

output "webook_urls" {
  description = "URLs of webhooks"
  value       = data.circleci_webhooks.project_webhooks.webhooks[*].url
}
