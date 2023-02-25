resource "circleci_checkout_key" "deploy_key" {
  project_slug = "github/acmeorg/foobar"
  type         = "deploy-key"
}

output "deploy_key_fingerprint" {
  description = "Fingerprint of SSH key"
  value       = circleci_checkout_key.deploy_key.fingerprint
}
