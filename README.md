Official repository and project is on [soucehut](https://git.sr.ht/~bacardi55/cleanABag). This repository is archived.

# cleanABag - A CLI tool to clean your (walla)bag

[![builds.sr.ht status](https://builds.sr.ht/~bacardi55/walgot.svg)](https://builds.sr.ht/~bacardi55/cleanABag?)
[![license: AGPL-3.0-only](https://img.shields.io/badge/license-AGPL--3.0--only-informational.svg)](LICENSE)

Official repository and project is on [sourcehut](https://git.sr.ht/~bacardi55/cleanABag). Github and codeberg are only mirrors. This is where binaries are uploaded.


cleanABag is a CLI tool for removing articles older than a given date from wallabag.

The goal is to avoid wasted storage of 100s or 1000s of articles that aren't needed in wallabag anymore.

To install, either download a binary from the [release page](https://git.sr.ht/~bacardi55/cleanABag/refs) or via the result of an [automated build](https://builds.sr.ht/?search=cleanABag).

Or manually (you need go >= 1.17 and make):

```bash
git clone https://git.sr.ht/~bacardi55/cleanABag
cd cleanABag
make dependencies
make build
```

The binary will be in the `bin` directory.



Usage


``` bash
cleanABag --help
```

``` text
Usage:
  cleanABag [command]

Available Commands:
  help        Help about any command
  prune       Delete old article from wallabag
  version     Print the version number of cleanABag

Flags:
  -c, --credentials string   config file with credentials to connect to wallabag (default is $HOME/.config/cleanABag/credentials.json) - Full path is needed.
  -h, --help                 help for cleanABag
  -v, --verbose              Verbose mode
```


``` bash
cleanABag prune --help
```
``` text
Delete old article from wallabag

Usage:
  cleanABag prune [flags]

Flags:
  -d, --date string   Articles older than this date will be removed if they match the archived/starred flags, format "YYYY-MM-DD".
      --delete        Delete articles. Without this flag, it will only do a dry run.
  -h, --help          help for prune
  -s, --starred       Include starred entry in deletion. False will prevent starred article to be deleted.
  -u, --unread        Include unread entries for deletion. False will prevent unread articles from being deleted
Global Flags:
  -c, --credentials string   config file with credentials to connect to wallabag (default is $HOME/.config/cleanABag/credentials.json) - Full path is needed.  -v, --verbose              Verbose mode
```


Example:

``` bash
# Remove archived articles older than 2021-12-31 (date in YYYY-MM-DD format) and that are not starred:
# Dry run:
cleanABag -c /path/to/credentials.json -d "2021-12-31"
# Delete for real:
cleanABag -c /path/to/credentials.json -d "2021-12-31" --delete

# Remove articles older than 2021-12-31, including unread but keep starred article
# Dry run:
cleanABag -c /path/to/credentials.json -d "2021-12-31" -u
# Delete for real:
cleanABag -c /path/to/credentials.json -d "2021-12-31" -u --delete

# Remove articles older than 2021-12-31, including unread and starred article
# Dry run:
cleanABag -c /path/to/credentials.json -d "2021-12-31" -u -s
# Delete for real:
cleanABag -c /path/to/credentials.json -d "2021-12-31" -u -s --delete

```


Example of `credentials.json`:

``` json
{
  "WallabagURL": "https://your.wallabag.tld",
  "ClientId": "client ID generate in your profile on wallabag"
  "ClientSecret": "client secrete generate in your profile on wallabag"
  "UserName": "your username",
  "UserPassword": "your password"
}
```
