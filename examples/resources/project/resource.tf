locals {
  project_slug = "github/acme/foobar"
}

resource "circleci_project" "my_project" {
  slug = local.project_slug
}

output "vcs_url" {
  description = "VCS url"
  value       = circleci_project.my_project.vcs_url
}
