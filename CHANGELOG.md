<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry is required to include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The tag should consist of where the change is being made ex. (x/staking), (store)
The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
"CLI Breaking" for breaking CLI commands.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased v1.1.0]

### Features

- (feat) [#51](https://github.com/IntelliXLabs/intellix/pull/51 ): feat: centralize the task dispather and merge the dispather and gateway into one process and start.
- (refactor) [#53](https://github.com/IntelliXLabs/intellix/pull/53 ): upgrade interactor to handle OperatorSocketUpdate event.

### Improvements

- (test) [#49](https://github.com/IntelliXLabs/intellix/pull/49) security: Use docker secrets to pass github token to prevent it from being leaked.
- (test) [#50](https://github.com/IntelliXLabs/intellix/pull/50 ): e2e: correct interactor config for e2e tests
- (feat) [#51](https://github.com/IntelliXLabs/intellix/pull/51 ): feat: centralize the task dispather and merge the dispather and gateway into one process and start.
- (datasource) [#52](https://github.com/IntelliXLabs/intellix/pull/52 ): feat: add additional data sources

## [Released v1.0.0]

After some time in development, we are now releasing version v1.0.0. This version features a price oracle that works with DVS. It also includes the launch of the public testnet.  

### Features

### Improvements

(test) [#39](https://github.com/IntelliXLabs/intellix/pull/39) improve: support price data source symbol converter.
(app) [#43](https://github.com/IntelliXLabs/intellix/pull/43) feat: add v1 upgrade handler, changelog version file, and changelog check in CI 
(test) [#45](https://github.com/IntelliXLabs/intellix/pull/45) improve: setup default `NO_PROXY` for `operator` container to avoid proxy issues.
(app) [#46](https://github.com/IntelliXLabs/intellix/pull/46) feat: update default denom to uitlx  
(app) [#47](https://github.com/IntelliXLabs/intellix/pull/47) chore: change denom from ITLX to IXN

### Bug Fixes
