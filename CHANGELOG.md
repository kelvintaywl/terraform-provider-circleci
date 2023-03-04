# Changelog

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
