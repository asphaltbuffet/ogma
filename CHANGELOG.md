# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Packaging now includes docs directory:
  - _README.md_
  - _LICENSE_
  - _CHANGELOG.md_
- Taskfile functionality added for local development

### Changed

- version information updated [#38](https://github.com/asphaltbuffet/ogma/issues/38)
- significant internal refactoring
- mocks are now automatically generated

### Removed

- completion command disabled [#37](https://github.com/asphaltbuffet/ogma/issues/37)
- docker images no longer generated as part of release artifacts
- add command has been removed in favor of using import

### Fixes

- base root command no longer returns an error [#35](https://github.com/asphaltbuffet/ogma/issues/35)
- default config file now included as part of packaging [#36](https://github.com/asphaltbuffet/ogma/issues/36)

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

[Unreleased]: https://github.com/asphaltbuffet/ogma/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/asphaltbuffet/ogma/releases/tag/v1.0.0
[0.0.2]: https://github.com/asphaltbuffet/ogma/releases/tag/v0.0.2
[0.0.1]: https://github.com/asphaltbuffet/ogma/releases/tag/v0.0.1

<!-- markdownlint-disable-file MD024 -->