# (Unofficial) Terraform Provider for CircleCI

## Support status

| Provider Block | Status |
| --- | --- |
| api_token | Done :white_check_mark: |
| hostname | Done :white_check_mark: |

| Data Source | Status |
| --- | --- |
| Webhook | In progress (10%) :construction_worker: |

| Resource | Status |
| --- | --- |
| Webhook | In progress (10%) :construction_worker: |

## Example

See [sandbox](sandbox/main.tf)

## Development

```console
# this project uses go 1.19
$ go mod tiny

# to build the go binary, and "install" to your local provider directory
# NOTE: this is a one-time action
$ make -f Makefile.dev init

# whenever we make changes to the code
$ make -f Makefile.dev build
# this then tries to terraform apply the sandbox
$ make -f Makefile.dev test_sandbox
```

## Notes

This depends on a CircleCI webhook Go SDK I (auto) generated:
https://github.com/kelvintaywl/circleci-webhook-go-sdk


## Why this is taking longer than I expected

1. I am unfortunately not a Go programmer. See [Hashicorp's stance on support for other languages here](https://developer.hashicorp.com/terraform/plugin/sdkv2/best-practices/other-languages)
2. There is [a template you are encouraged to use](https://github.com/hashicorp/terraform-provider-scaffolding-framework) but the internal implementation is missing, so not knowing how to use [the framework](https://github.com/hashicorp/terraform-plugin-framework) and Go makes it tougher.
3. The provided [framework (SDK)](https://github.com/hashicorp/terraform-plugin-framework) is new, so older providers are not using it, which also makes it difficult to refer to example codes.

