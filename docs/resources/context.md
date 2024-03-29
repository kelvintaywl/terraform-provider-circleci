---
page_title: "circleci_context Resource - terraform-provider-circleci"
subcategory: ""
description: |-
  Manages a context
---

# circleci_context (Resource)

Manages a context

**Note**: Contexts cannot be updated.

If you modify the `name`, this will delete the existing context and recreate one instead.

Creating a context for an account is only supported in CircleCI Server (self-hosted) instance.
A _warning_ message is printed to remind users about this, if owner type is selected as `account`.
You can observe this message by varying the `TF_LOG` value.

## Example Usage

```terraform
resource "circleci_context" "from_tf" {
  name = "from_tf"
  owner = {
    // replace id with your organization ID
    id   = "7f284df8-ac74-42d5-9fad-ab23f731e475"
    type = "organization"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the context
- `owner` (Attributes) The owner of the context (see [below for nested schema](#nestedatt--owner))

### Read-Only

- `created_at` (String) The date and time the schedule was created
- `id` (String) The unique ID of the context

<a id="nestedatt--owner"></a>
### Nested Schema for `owner`

Required:

- `id` (String) The unique ID of the owner
- `type` (String) The type of the owner. Accepts `account` or `organization`. Accounts are only used as context owners in **Server**.

## Import

An existing context can be imported via its owner type, owner ID and unique context ID (UUID).

```console
# import an organization context
$ terraform import circleci_context.my_context "organization,<ORG ID>,<CONTEXT ID>"
```
