# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Delete command now available to remove datastore
- Export command now available for mail/listings
  - Importing this file will have no effect unless changes are made to the json file. Then it serves as an editing method.

### Changed

- Taskfile now has `snapshot` that replaces previous `build` task
- `build` task now uses `go build` to compile a local binary for dev use
- Additional linting rules and settings update
- Import command now autoimports listings or mail records.

### Fixes

- Fixed potential panic areas in unit tests where string length could go out of bounds

## [1.1.1] - 2021-12-22

### Added

- Available via snap (candidate-only)
  - known issue with config file not working

### Changed

- Help text for commands has been updated

### Fixes

- Unit tests updated to pass on windows
- File closing fixed in multiple locations

## [1.1.0] - 2021-12-19

### Added

- Add mail subcommand for tracking correspondence
- Packaging now includes docs directory:
  - _README.md_
  - _LICENSE_
  - _CHANGELOG.md_
- Taskfile functionality added for local development

### Changed

- version information updated [#38](https://github.com/asphaltbuffet/ogma/issues/38)
- significant internal refactoring
- mocks are now automatically generated
- search command returns correspondence information
- import now has subcommands for mail or listing sources
- syslogging only available for Darwin and Linux builds
- datastore testing better isolated from rest of application

### Removed

- completion command disabled [#37](https://github.com/asphaltbuffet/ogma/issues/37)
- docker images no longer generated as part of release artifacts
- add command has been removed in favor of using import

### Fixes

- base root command no longer returns an error [#35](https://github.com/asphaltbuffet/ogma/issues/35)
- default config file now included as part of packaging [#36](https://github.com/asphaltbuffet/ogma/issues/36)
- json files are properly closed after being processed
- searching does not create a blank datastore if one didn't exist before

## [1.0.0] - 2021-11-15

### Added

- Initial functional release
- Importing listings (one-to-many records)
- Adding listings (single record)
- Search (by member number only)
- Application configuration
- Usage documentation

### Changed

- Update to go 1.17.3

## [0.0.2] - 2021-11-01

### Added

- Application configuration file (.ogma)
- Logging
- New commands
  - _Issues_
  - _Listings_
    - _add_
    - _search_
- New _info_ flag for base command

### Changed

- More unit testing

## [0.0.1] - 2021-10-21

### Added

- Initial commit [Go Repository Template](https://github.com/golang-templates/seed)

### Changed

### Deprecated

### Removed

### Fixes

### Security

[Unreleased]: https://github.com/asphaltbuffet/ogma/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/asphaltbuffet/ogma/releases/tag/v1.2.0
[1.1.1]: https://github.com/asphaltbuffet/ogma/releases/tag/v1.1.1
[1.1.0]: https://github.com/asphaltbuffet/ogma/releases/tag/v1.1.0
[1.0.0]: https://github.com/asphaltbuffet/ogma/releases/tag/v1.0.0
[0.0.2]: https://github.com/asphaltbuffet/ogma/releases/tag/v0.0.2
[0.0.1]: https://github.com/asphaltbuffet/ogma/releases/tag/v0.0.1

<!-- markdownlint-disable-file MD024 -->
