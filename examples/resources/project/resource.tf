# set up a new CircleCI project
# ASSUMPTION: the GitHub project has a .circleci/config.yml on its default branch
resource "circleci_project" "my_project" {
  slug = "github/acme/foobar"
}

# add a project env var to this project
resource "circleci_env_var" "my_env_var" {
  project_slug = circleci_project.my_project.slug
  name         = "FOOBAR"
  value        = "0Cme2FmlXk"
}

output "vcs_url" {
  description = "VCS url"
  value       = circleci_project.my_project.vcs_url
}
