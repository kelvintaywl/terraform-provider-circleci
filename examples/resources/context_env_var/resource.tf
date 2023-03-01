locals {
  // replace with your organization ID
  org_id = "7f284df8-ac74-42d5-9fad-ab23f731e475"
}

resource "circleci_context" "example" {
  name = "example"
  owner = {
    id   = local.org_id
    type = "organization"
  }
}

resource "circleci_context_env_var" "test_envvar" {

  for_each = {
    FOOBAR   = "Lorem Ipsum"
    FIZZBUZZ = "random1234"
  }

  name       = each.key
  value      = each.value
  context_id = circleci_context.example.id
}
