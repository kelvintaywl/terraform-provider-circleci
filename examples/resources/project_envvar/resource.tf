locals {
  // Replace this with your CircleCI project slug
  project_slug = "github/acmeorg/foobar"
}

resource "circleci_project_envvar" "project_envvar_foobar" {
  project_slug = local.project_slug
  name         = "FOOBAR"
  value        = "0Cme2FmlXk"
}

resource "circleci_project_envvar" "project_envvar_foobar_fizzbuzz" {
  project_slug = local.project_slug
  name         = "FIZZBUZZ"
  value        = "Vbt2efixZAkrmTYiirhd"
}
