data "circleci_runner_tokens" "test" {
  resource_class = "kelvintaywl-tf/test"
}

output "tokens" {
  description = "runner tokens"
  value       = data.circleci_runner_tokens.test.tokens
}