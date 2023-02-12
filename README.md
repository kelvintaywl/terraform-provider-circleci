# (Unofficial) Terraform Provider for CircleCI

[![Go Report Card](https://goreportcard.com/badge/github.com/kelvintaywl/terraform-provider-circleci)](https://goreportcard.com/report/github.com/kelvintaywl/terraform-provider-circleci)
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kelvintaywl-cci/terraform-provider-circleci/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/kelvintaywl-cci/terraform-provider-circleci/tree/main)

## Support status

| Provider Block | Status | Remarks |
| --- | --- | --- |
| api_token | Done :white_check_mark: | $CIRCLE_TOKEN supported |
| hostname | Done :white_check_mark: | $CIRCLE_HOSTNAME supported |

| Data Source | Status | Remarks |
| --- | --- | --- |
| Webhooks | In progress (90%) :construction_worker: | TODO: pagination |

| Resource | Status | Remarks |
| --- | --- | --- |
| Webhook | In progress (90%) :construction_worker: | TODO: support importing of state |

## Example

See [sandbox](sandbox/main.tf)

## Development

```console
# this project uses go 1.19
$ go mod tiny

# to build the go binary, and "install" to your local provider directory
# NOTE: this is a one-time action
$ make init

# whenever we make changes to the code
$ make build
# this then tries to terraform apply the sandbox
$ make tf.plan

# when you are ready to apply
$ make tf.apply

# destroy, as needed
$ make tf.destroy
```

## Testing

This uses a dummy CircleCI project for acceptance tests:
https://github.com/kelvintaywl-cci/tf-provider-acceptance-test-dummy

```console
# Run acceptance tests
$ export CIRCLE_TOKEN="user API token that can CRUD the dummy project"
$ make testacc
```


## Docs

```console
# to generate docs
# check docs/index.md
$ make docs
```

See [docs](docs/index.md)


## Notes

This depends on a CircleCI webhook Go SDK I (auto) generated:
https://github.com/kelvintaywl/circleci-webhook-go-sdk


## Why this is taking longer than I expected

1. I am unfortunately not a Go programmer. See [Hashicorp's stance on support for other languages here](https://developer.hashicorp.com/terraform/plugin/sdkv2/best-practices/other-languages)
2. There is [a template you are encouraged to use](https://github.com/hashicorp/terraform-provider-scaffolding-framework) but the internal implementation is missing, so not knowing how to use [the framework](https://github.com/hashicorp/terraform-plugin-framework) and Go makes it tougher.
3. The provided [framework (SDK)](https://github.com/hashicorp/terraform-plugin-framework) is new, so older providers are not using it, which also makes it difficult to refer to example codes.
