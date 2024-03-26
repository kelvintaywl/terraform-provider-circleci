resource "circleci_env_var" "my_env_vars" {
  for_each = {
    "FOOBAR"   = "0Cme2FmlXk"
    "FIZZBUZZ" = "Vbt2efixZAkrmTYiirhd"
  }

  project_slug = "circleci/7UQdtYSr1caLbAR2cHJdU7/2DACeEvUr7MosidActmnUs"
  name         = each.key
  value        = each.value
}
