locals {
  // Replace this with your CircleCI project ID
  project_id = "c124cca6-d03e-4733-b84d-32b02347b78c"
}


data "circleci_webhooks" "webhooks" {
  project_id = local.project_id
}

output "webhooks" {
  description = "current webhooks"
  value       = data.circleci_webhooks.webhooks.webhooks[*]
  // required since signing_secret is sensitive
  sensitive = true
}
