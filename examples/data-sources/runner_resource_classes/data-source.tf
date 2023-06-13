data "circleci_runner_resource_classes" "tf" {
  namespace = "kelvintaywl-tf"
}

output "resources" {
  description = "runner resource-classes"
  value       = data.circleci_runner_resource_classes.tf.resource_classes
}
