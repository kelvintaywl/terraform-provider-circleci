data "circleci_checkout_keys" "my_keys" {
  project_slug = "github/acmeorg/foobar"
}

output "keys" {
  description = "all keys for this project"
  value       = data.circleci_checkout_keys.my_keys.keys
}
