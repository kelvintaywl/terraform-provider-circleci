locals {
  namespace      = "kelvintaywl-tf"
  resource_class = "from-tf"
}

resource "circleci_runner_resource_class" "from_tf" {
  resource_class = "${local.namespace}/${local.resource_class}"
  description    = "Test from Terraform"
}

output "runner_from_tf_id" {
  description = "runner resource class ID"
  value       = circleci_runner_resource_class.from_tf.id
}
