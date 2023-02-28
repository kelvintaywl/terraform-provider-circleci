data "circleci_context" "test" {
  name = "my_context"
  owner = {
    id   = "7f284df8-ac74-42d5-9fad-ab23f731e475"
    type = "organization"
  }
}

output "context_id" {
  description = "my context id"
  value       = data.circleci_context.test.id
}
