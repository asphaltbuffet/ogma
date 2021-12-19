# OGMA

A LEX Magazine DB and letter tracking application.

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![GitHub Release](https://img.shields.io/github/v/release/asphaltbuffet/ogma)](https://github.com/asphaltbuffet/ogma/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/asphaltbuffet/ogma.svg)](https://pkg.go.dev/github.com/asphaltbuffet/ogma)
[![go.mod](https://img.shields.io/github/go-mod/go-version/asphaltbuffet/ogma)](go.mod)
[![LICENSE](https://img.shields.io/github/license/asphaltbuffet/ogma)](LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/asphaltbuffet/ogma/build)](https://github.com/asphaltbuffet/ogma/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/asphaltbuffet/ogma)](https://goreportcard.com/report/github.com/asphaltbuffet/ogma)
[![Codecov](https://codecov.io/gh/asphaltbuffet/ogma/branch/main/graph/badge.svg)](https://codecov.io/gh/asphaltbuffet/ogma)

- [OGMA](#ogma)
  - [Usage](#usage)
    - [Mail Command](#mail-command)
    - [Import Command](#import-command)
      - [Example import file](#example-import-file)
    - [Search Command](#search-command)
  - [Configuration](#configuration)
    - [Default config](#default-config)

## Usage

Ogma does not come pre-loaded with any issue information. As the service requires members to buy issues but not pay for letters forwarded, this application is not intended to skirt that monetary flow. You can hand-type information, scan, invoke blood magic, whatever.

There is one method to add issue listings for later use: [Import](#import-command).

Application usage details can be found via the `-h` or `--help` flag in the base application or with any command. This will show argument details, flag detail, and some examples. In-application documentation always supercedes this documentation.

### Mail Command

The mail command tracks correspondence and provides the ability to link with LEX listings, members, and/or other correspondence.

```bash
ogma mail --sender=<member> --receiver=<member> --date=<date> -link<member or mail ref>
```

By default, the member numbers are set to the configured member number in the application configuration. They must be entered as integers. Any extended member number (those with a letter) should be entered without any text characters.

On success, a reference number is returned for use in tracking correspondence artifacts.

```bash
ogma mail -s1234 -r5678 -d2021-11-15
Added mail. Reference: f8427e
```

### Import Command

The import command takes the filename (for now) of a json file that contains listing or mail entries.

**There is no checking for duplicates already in the application database. Careful!**

```bash
ogma import [listing|mail] <filename.json>
```

By default, the import command will only output the number of entries saved to the db ([#33](https://github.com/asphaltbuffet/ogma/issues/33)). No, there's no way to check this right now other than doing manual searches for entries to figure out what made it in.

If you want to see everything that has been imported, use the verbose flag (`-v` or `--verbose`) to see all entries printed to screen. This may be a lot of stuff on your screen...

```bash
ogma import listing <filename> -v
```

#### Example listing import file

```json
{
    "listings": [
        {
            "volume": 1,
            "issue": 1,
            "year": 1986,
            "season": "Spring",
            "page": 1,
            "category": "Art & Photography",
            "member": 123,
            "alt": "",
            "international": true,
            "review": false,
            "text": "This is an example.",
            "art": false,
            "flag": false
        }
    ]
}
```

### Search Command

This is the primary use of the application and is simplified at the moment. Only searching by member number is supported, and this must be entered as an integer. Adding in a letter after the member number as seen in some issues is invalid and will fail. All listings with that member number will be found (listings with an alphabetic extension will be included and shown as such).

```bash
ogma search <member number>
```

Currently the output defaults to something pretty, with colors (results on windows may vary). Output configuration may be expanded in the future. See [#16](https://github.com/asphaltbuffet/ogma/issues/16)

## Configuration

### Default config

```yaml
logging:
  level: info
search:
  max_results: 10
datastore:
  filename: "ogma.db"
defaults:
  issue: 56
  max_column: 40
member: 13401
```
