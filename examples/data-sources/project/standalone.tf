data "circleci_project" "test" {
  slug = "circleci/7UQdtYSr1caLbAR2cHJdU7/2DACeEvUr7MosidActmnUs"
}

output "url" {
  description = "project_url"
  value       = data.circleci_project.test.vcs_url
}