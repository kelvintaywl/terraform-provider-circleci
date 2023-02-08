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

data "circleci_webhook" "project_webhook" {
  project_id = "32bbe47f-2bdf-4bb7-8390-ce682161a95f"
}
