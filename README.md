# OGMA
A LEX Magazine DB ~~and letter tracking application~~(coming soon).

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![GitHub Release](https://img.shields.io/github/v/release/asphaltbuffet/ogma)](https://github.com/asphaltbuffet/ogma/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/asphaltbuffet/ogma.svg)](https://pkg.go.dev/github.com/asphaltbuffet/ogma)
[![go.mod](https://img.shields.io/github/go-mod/go-version/asphaltbuffet/ogma)](go.mod)
[![LICENSE](https://img.shields.io/github/license/asphaltbuffet/ogma)](LICENSE)
[![Build Status](https://img.shields.io/github/workflow/status/asphaltbuffet/ogma/build)](https://github.com/asphaltbuffet/ogma/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/asphaltbuffet/ogma)](https://goreportcard.com/report/github.com/asphaltbuffet/ogma)
[![Codecov](https://codecov.io/gh/asphaltbuffet/ogma/branch/main/graph/badge.svg)](https://codecov.io/gh/asphaltbuffet/ogma)

## Usage
Ogma does not come pre-loaded with any issue information. As the service requires members to buy issues but not pay for letters forwarded, this application is not intended to skirt that monetary flow. You can hand-type information, scan, invoke blood magic, whatever.

There are two methods to add issue listings for later use, _Add_ and _Import_.
### Add Command
The add command allows single entries to be saved. It is very manual and is really only good for occasional data entry. Seriously, use import.
```
ogma add [FLAGS]
```
Ok, fine. Each listing field is a separate flag. Not all flags are required, but the needed ones are.

TODO: details on flags and fields go here

### Import Command
The import command takes the filename (for now) of a json file that contains listing entries. 

**There is no checking for duplicates already in the application database. Careful!**
```
ogma import <filename>
```

By default, the import command will only output the number of listings saved to the db (#33). No, there's no way to check this right now other than doing manual searches for entries to figure out what made it in.

If you want to see everything that has been imported, use the verbose flag (`-v` or `--verbose`) to see all entries printed to screen. This may be a lot of stuff on your screen...

```
ogma import <filename> -v
```

#### Example import file
TODO: put the example here...

### Search Command
This is the primary use of the application and is simplified at the moment. Only searching by member number is supported, and this must be entered as an integer. Adding in a letter after the member number as seen in some issues is invalid and will fail. All listings with that member number will be found (listings with an alphabetic extension will be included and shown as such).

```
ogma search <member number>
```

Currently the output defaults to something pretty, with colors. This may be enhanced in the future. #16

## Configuration
### Default config (.ogma)
```yaml
logging:
    level: warn
search:
    max_results: 10
datastore:
    filename: "ogma.db"
defaults:
    issue: 56
    max_column: 40
```