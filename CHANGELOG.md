# Changelog

## [1.0.0] - 2024-03-03

### Updated

- Update docs around runner resource class & example usage
- [Breaking] Update license to Mozilla Public License Version 2.0

## [0.12.0] - 2023-08-29

### Added

- Support project resource

## [0.11.1] - 2023-08-12

### Added

- Support imports for Context and Runner resource-class resource

## [0.11.0] - 2023-08-05

### Fixed

- Update Webhook resource & data-source to detect changes on events and signing_secret as expected

## [0.10.4] - 2023-06-28

### Updated

- Update Runner Token resource examples

## [0.10.3] - 2023-06-15

### Fixed

- Use server hostname as-is for Runner API calls

## [0.10.2] - 2023-06-15

### Updated

- Update main doc to latest list of resources supported

## [0.10.1] - 2023-06-15

### Fixed

- Remove duplicate info in document

## [0.10.0] - 2023-06-15

### Added

- Support runner resource-class as resource
- Support runner resource-classes as data source
- Support runner token as resource
- Support runner tokens as data source

## [0.9.3] - 2023-03-05

### Fixed

- Skip NotFound errors when destroying resource, as [recommended](https://developer.hashicorp.com/terraform/plugin/framework/resources/delete#recommendations)

## [0.9.2] - 2023-03-05

### Changed

- Document maximum no. of webhooks per project (data source)
- Implement pagination for checkout keys (data source)

## [0.9.1] - 2023-03-04

### Changed

- Refactored warning messages for call-to-actions (improved developer experience)

## [0.9.0] - 2023-03-01

### Added

- Support context environment variable as resource

## [0.8.2] - 2023-02-28

### Added

- Support context as data source

## [0.8.1] - 2023-02-27

### Changed

- Corrected documentation around context

## [0.8.0] - 2023-02-27

### Added

- Support context as resource

## [0.7.2] - 2023-02-26

### Added

- Support import functionality for Schedule and Webhook resources

## [0.7.1] - 2023-02-26

### Added

- Log warning about users needed to delete public key of checkout keys in VCS

### Changed

- Refactored handling of unchanged attributes when updating
- Refactored handling of resources that do not support updates (env_var, checkout_key)

## [0.7.0] - 2023-02-25

### Added

- Support checkout keys as data source
- Support checkout key as resource

## [0.6.0] - 2023-02-25

### Changed

- **BREAKING** Refactored project environment variable (resource) to `env_var`

## [0.5.0] - 2023-02-25

### Added

- Support project as data source

## [0.4.0] - 2023-02-24

### Added

- Support project environment variable as resource

## [0.3.1] - 2023-02-21

### Added

- Updated provider documentation

## [0.3.0] - 2023-02-21

### Added

- Support scheduled pipeline (schedule) as resource

## [0.2.1] - 2023-02-14

### Fixed

- Removed unnecessary print logs

## [0.2.0] - 2023-02-12

### Fixed

- Fixed provider path

## [0.1.0] - 2023-02-12

### Added

- First release with webhook resource & webhooks data source
