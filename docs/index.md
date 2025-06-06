---
page_title: "Provider: CircleCI"
description: |-
  The CircleCI provider enpowers users to CircleCI resources via Terraform
---

# CIRCLECI Provider

**This is an unofficial Terraform provider for CircleCI.**

The CircleCI provider supports the creation of specific CircleCI resources.

Currently, the following resources are supported:

- [Project](https://circleci.com/docs/create-project/)
- [Project Webhook](https://circleci.com/docs/webhooks/)
- [Project Scheduled Pipeline](https://circleci.com/docs/scheduled-pipelines/)
- [Project Environment Variables](https://circleci.com/docs/set-environment-variable/#set-an-environment-variable-in-a-project)
- Project Checkout key
- [Context](https://circleci.com/docs/contexts/)
- [Context Environment Variables](https://circleci.com/docs/contexts/#adding-and-removing-environment-variables-from-restricted-contexts)
- [Runner Resource-class](https://circleci.com/docs/runner-faqs/#what-is-a-runner-resource-class)
- [Runner Token](https://circleci.com/docs/runner-faqs/#what-is-a-runner-resource-class)

Use the navigation to the left to read about the available resources.

## Example Usage

```terraform
terraform {
  required_providers {
    circleci = {
      source = "kelvintaywl/circleci"
      # update version accordingly
      version = "0.12.0"
    }
  }
}

provider "circleci" {
  // You can also set this via CIRCLE_TOKEN environment variable.
  api_token = "myCircleCIUserAPIToken"

  // Defaults to circleci.com
  // If you are using a self-hosted CircleCI instance (aka Server),
  // specify your self-hosted server's domain here ('https://' not required).
  // This can also be set via CIRCLE_HOSTNAME environment variable,
  hostname = "circleci.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_token` (String) A CircleCI user API token. This can also be set via the `CIRCLE_TOKEN` environment variable.
- `hostname` (String) CircleCI hostname (default: circleci.com). This can also be set via the `CIRCLE_HOSTNAME` environment variable.
- `max_retries` (Number) Maximum number of retries for API calls when retry is enabled (default: 3).
- `retry` (Boolean) Whether to retry API calls when provider receives an HTTP 429 status code (default: false).
