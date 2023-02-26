# (Unofficial) Terraform Provider for CircleCI

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kelvintaywl-cci/terraform-provider-circleci/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/kelvintaywl-cci/terraform-provider-circleci/tree/main)

https://registry.terraform.io/providers/kelvintaywl/circleci/latest

## Support status

| Provider Block | Status | Remarks |
| --- | --- | --- |
| api_token | Done :white_check_mark: | $CIRCLE_TOKEN supported |
| hostname | Done :white_check_mark: | $CIRCLE_HOSTNAME supported |

| Data Source | Status | Remarks |
| --- | --- | --- |
| Webhooks | In progress (90%) :construction_worker: | TODO: pagination |
| Project | Done :white_check_mark: | |
| Checkout keys | Done :white_check_mark: | |

| Resource | Status | Remarks |
| --- | --- | --- |
| Webhook | In progress (90%) :construction_worker: | TODO: support importing of state |
| Schedule | In progress (90%) :construction_worker: | TODO: support importing of state |
| Project Environment Variable | Done :white_check_mark: | |
| Checkout key | Done :white_check_mark: | |

## Example

See [sandbox](sandbox/main.tf)

## Development

```console
# this project uses go 1.19
$ go mod download

# or, if you want to upgrade the Go dependencies too
$ go get -u
$ go mod tidy

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

## Releasing

Currently, we release locally with [GoReleaser](https://goreleaser.com/install/).
The config can be found in [.goreleaser.yml](.goreleaser.yml)

```console
# NOTE: make sure documents are generated!

# tag release
$ git tag vX.Y.Z -m "some message"

# install GoReleaser on MacOS if required
# See https://goreleaser.com/install/
$ brew install goreleaser/tap/goreleaser

# set required env vars
$ export GITHUB_TOKEN="your GitHub Token value, with public_repo scope required"
$ export GPG_FINGERPRINT="your GPG fingerprint, registered to your Terraform namespace"
# optional
$ export GPG_TTY=$(tty)

# make a GitHub release (with artifacts),
# via GoReleaser
$ goreleaser release --clean
```

## Notes

This depends on a CircleCI Go SDK I (auto) generated:
https://github.com/kelvintaywl/circleci-go-sdk


## Why this is taking longer than I expected

1. I am unfortunately not a Go programmer. See [Hashicorp's stance on support for other languages here](https://developer.hashicorp.com/terraform/plugin/sdkv2/best-practices/other-languages)
2. There is [a template you are encouraged to use](https://github.com/hashicorp/terraform-provider-scaffolding-framework) but the internal implementation is missing, so not knowing how to use [the framework](https://github.com/hashicorp/terraform-plugin-framework) and Go makes it tougher.
3. The provided [framework (SDK)](https://github.com/hashicorp/terraform-plugin-framework) is new, so older providers are not using it, which also makes it difficult to refer to example codes.
