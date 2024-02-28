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
