# using for_each to generate multiple env var resources
resource "circleci_env_var" "my_env_vars" {
  for_each = {
    "FOOBAR"   = "0Cme2FmlXk"
    "FIZZBUZZ" = "Vbt2efixZAkrmTYiirhd"
  }

  project_slug = "github/acmeorg/foobar"
  name         = each.key
  value        = each.value
}
