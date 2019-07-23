# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- `config` command to get and set configuration options.
- `endpoint list|ls` command to print the endpoints list as a table.
  - `--format` flag to format output using a Go template.
- `stack list|ls` command to print the stacks list as a table.
  - `--swarm` flag sets a filter by swarm Id.
  - `--endpoint` flag sets a filter by endpoint Id.
  -`-q, --quiet`) flag causes stack list to print only stack names.
  - `--format` flag to format output using a Go template.
- `stack deploy|up|create` command to deploy a stack.
  - `-c, --stack-file` flag sets the stack file to use.
  - `-e, --env-file` flag sets the environment variables file to use.
  - `--replace-env` flag causes environment variables to be replaced instead of merged while updating a stack.
  - `--endpoint` flag sets the endpoint to use.
  - `-p, --prune` flag causes services that are no longer referenced to be removed.
- `stack remove|rm|down` command to remove a stack.
  - `--strict` flag causes a failure if the stack does not exist.
- `status` command to show Portainer status as a table.
  - `--format` flag to format output using a Go template.
- `help` command and `-h, --help` global flag to print global help.
- `-h, --help` flags on each command to print its help.
- `-t, --timeout` global flag to set a timeout for requests execution.
- `--config` global flag to set the path to a configuration file.
- `--version` global flag to print the client version.
- `completion` command to print Bash completion script.

### Removed
- `-a` flag, which used to select an action to execute. There is a command now for each action.
- `bash`, `jq` and `httpie` programs in the Docker image. The client doesn't use them anymore.

### Changed
- Single executable binary instead of a bash script. Project has been rewritten in Go language.
- No external programs (bash, httpie, jq) dependency. The new executable binary is self-sufficient.
- `-u` global flag renamed to `--user`.
- `-p` global flag renamed to `--password`.
- `-l` global flag renamed to `--url`.
- `-s` global flag renamed to `--insecure`. It does not receive a value anymore, it is a boolean flag.
- `-v` global flag renamed to `-v, --verbose`. It does not receive a value anymore, it is a boolean flag.
- `-d` global flag renamed to `-d, --debug`. It does not receive a value anymore, it is a boolean flag.
- `-n` global flag moved as a parameter to the `stack deploy` command.
- `-c` global flag renamed to `-c, --stack-file` and moved to the `stack deploy` command.
- `-g` global flag renamed to `-e, --env-file` and moved to the `stack deploy` command.
- `-r` global flag renamed to `-p, --prune` and moved to the `stack deploy` command. It does not receive a value anymore, it is a boolean flag.
- `-e` global flag renamed to `--endpoint` and moved to the `stack deploy` command.
- `-t` global flag renamed to `--strict` and moved to the `stack remove` command.
- All environment variables prefixed with "PSU_" and renamed to match command and flag names.

## [0.1.1] - 2019-06-05
### Fixed
- Fixed error when environment variables loaded from file contain spaces in their values [#14](https://github.com/greenled/portainer-stack-utils/pull/14)

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

[Unreleased]: https://github.com/greenled/portainer-stack-utils/compare/0.1.1...HEAD
[0.1.1]: https://github.com/greenled/portainer-stack-utils/releases/tag/0.1.1
[0.1.0]: https://github.com/greenled/portainer-stack-utils/releases/tag/0.1.0
