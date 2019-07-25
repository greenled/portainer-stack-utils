# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2019-07-25
### Added
- New actions: `ls`, `status`, `services`, `tasks`, `tasks:healthy`, `containers`, `login`, `lint`, `inspect`, `system:info`, `actions`, `help` and `version`
- New options: `--auth-token=[AUTH_TOKEN]`,	`--compose-file-base64=[BASE64]`, `--env-file-base64=[BASE64]`, `--timeout=[SECONDS]`, `--detect-job=[true|false]`, `--service=[SERVICE_NAME]`, `--insecure`, `--masked-variables`, `--quiet`, `--lint`, `--help` and `--version`
- New flags: `-A`, `-C`, `-F`, `-G`, `-T`, `-j`, `-i`, `-S`, `-m`, `-q`, `-L`, `-h` and `-V`
- New environment variables: `PORTAINER_AUTH_TOKEN`, `TIMEOUT`, `AUTO_DETECT_JOB`, `PORTAINER_SERVICE_NAME`, `MASKED_VARIABLES`, `QUIET_MODE` and `DOCKER_COMPOSE_LINT`
- The Docker image include now `docker-compose` to be able to lint Docker compose/stack file
- The `core` Docker image variant doesn't include `docker-compose`, so it's a bit smaller. But you can't lint Docker compose/stack file before deploying a stack
- The `debian` and `debian-core` Docker image variants, use [Debian](https://www.debian.org) instead of [Alpine](https://alpinelinux.org/) as base image for `psu`
- Online documentation via [docsify](https://docsify.js.org)
- Tests who run automatically on each git push via [GitLab CI](https://docs.gitlab.com/ce/ci/)

### Changed
- The `undeploy` action is now an aliased action. You should use `rm` action instead

### Deprecated
- The `--secure=[yes|no]` option and `-s` flag are deprecated. Use the `--insecure` option instead (`psu <action> ... --inscure`)
- The `--action=[ACTION_NAME]` option and `-a` flag are deprecated. Use `<action>` argument instead (`psu <action> ...`)

## [0.1.1] - 2019-06-05
### Fixed
- Fixed error when environment variables loaded from file contain spaces in their values [#14](https://gitlab.com/psuapp/psu/merge_requests/14)

## [0.1.0] - 2019-05-24
### Added
- Stack deployment
- Stack update
- Stack undeployment
- Configuration through environment variables
- Configuration through flags
- Stack environment variables loading from file
- Optional SSL verification of Portainer instance
- Verbose mode
- Debug mode
- Strict mode

[Unreleased]: https://gitlab.com/psuapp/psu/compare/v0.1.1...master
[0.1.1]: https://gitlab.com/psuapp/psu/-/tags/v0.1.1
[0.1.0]: https://gitlab.com/psuapp/psu/-/tags/v0.1.0
