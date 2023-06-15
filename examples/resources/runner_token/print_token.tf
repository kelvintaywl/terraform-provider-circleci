resource "circleci_runner_token" "my_new_token" {
  resource_class = "kelvintaywl-tf/test"
  nickname       = "my-new-token"
}

output "token" {
  description = "kelvintaywl-tf/test runner token (my-new-token)"
  # to print token raw in the output
  value = nonsensitive(circleci_runner_token.my_new_token.token)
}
