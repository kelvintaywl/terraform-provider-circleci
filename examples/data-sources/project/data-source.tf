data "circleci_project" "test" {
  slug = "github/acmeorg/foobar"
}

output "url" {
  description = "project_url"
  value       = data.circleci_project.test.vcs_url
}