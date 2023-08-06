# Change Log
All notable changes to this project will be documented in this file.

## [v0.10.0](https://github.com/gliderlabs/registrator/compare/v0.9.0...v0.10.0) - 2023-08-06

### Added

- #82 @josegonzalez Add support for raspbian/bullseye on arm64 architectures
- #98 @Coffee2CodeNL: Add Bookworm Support

### Changed

- #81 @dependabot: chore(deps): bump golang from 1.18.1-buster to 1.18.2-buster
- #83 @josegonzalez Publish arm64 packages for debian instead of raspbian
- #94 @dependabot: chore(deps): bump golang from 1.18.2-buster to 1.20.1-buster
- #95 @dependabot: chore(deps): bump golang from 1.20.1-buster to 1.20.4-buster
- #96 @josegonzalez chore: update runner to use ubuntu 20.04
- #97 @dependabot: chore(deps): bump golang from 1.20.4-buster to 1.20.5-buster

## [v0.9.0](https://github.com/gliderlabs/registrator/compare/v0.8.1...v0.9.0) - 2022-05-10

### Added

- #74 @josegonzalez Publish armhf package to ubuntu/focal
- #78 @josegonzalez Publish package for Ubuntu 22.04

### Changed

- #73 @dependabot chore(deps): bump golang from 1.17.7-buster to 1.17.8-buster
- #76 @dependabot chore(deps): bump golang from 1.17.8-buster to 1.18.1-buster
- #77 @chibadaijiro Specify the full path in the module directive of go.mod.
- #79 @josegonzalez Update go modules

## [v0.8.1](https://github.com/gliderlabs/registrator/compare/v0.8.0...v0.8.1) - 2022-03-05

### Added

- #70: RyanGaudion Add Raspbian Bullseye

### Changed

- #69: @dependabot chore(deps): bump golang from 1.17.6-buster to 1.17.7-buster

## [v0.8.0](https://github.com/gliderlabs/registrator/compare/v0.7.1...v0.8.0) - 2022-01-31

### Added

- #66: @josegonzalez Add support for arm64

### Fixed

- #67: @josegonzalez Add missing jq dependency to build environment

### Changed

- #62: @dependabot chore(deps): bump golang from 1.17.2-buster to 1.17.3-buster
- #65: @dependabot chore(deps): bump golang from 1.17.3-buster to 1.17.6-buster

## [v0.7.1](https://github.com/gliderlabs/registrator/compare/v0.7.0...v0.7.1) - 2020-10-28

### Fixed

- #60: @josegonzalez Ensure binary name in package is sigil

## [v0.7.0](https://github.com/gliderlabs/registrator/compare/v0.6.0...v0.7.0) - 2020-10-28

### Fixed

- #56: @0xflotus Fix typo in readme

### Added

- #58: @josegonzalez Add arm support
- #57: @adam12 Add bullseye to deb release task

## [v0.6.0](https://github.com/gliderlabs/registrator/compare/v0.5.0...v0.6.0) - 2020-05-06

### Added

- #53: @josegonzalez Release packages for Focal Fossa

## [v0.5.0](https://github.com/gliderlabs/registrator/compare/v0.4.0...v0.5.0) - 2020-03-13
### Fixed

- #19: @stormcat24 Use correct protocol for download url
- #24: @mozamimy Fix a typo in SplitKv function
- #30: @lalyos Fix tojson
- #39: @josegonzalez Fix sigil installation instructions
- #46: @uphy `exists` function never returns false on full path

### Added

- #13: @lalyos Make go templating left and right delimiter customizable.
- #16: @lalyos Add substring function
- #21: @lalyos Add base64enc and base64dec functions
- #28: @lalyos Add [jmespath](http://jmespath.org) function
- #52: @josegonzalez Release packages in CI

### Changed

- @josegonzalez Upgrade to circleci 2.1
- @josegonzalez Upgrade to golang 1.13.8
