resource "circleci_context" "from_tf" {
  name = "from_tf"
  owner = {
    // replace id with your organization ID
    id   = "7f284df8-ac74-42d5-9fad-ab23f731e475"
    type = "organization"
  }
}
