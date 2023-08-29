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
  project_slug = "github/kelvintaywl-cci/delete-api"
}

resource "circleci_project" "my_project" {
  slug = local.project_slug
}

output "vcs_url" {
  description = "VCS url"
  value       = circleci_project.my_project.vcs_url
}
