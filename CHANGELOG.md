# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- All commands are self documented.
- If an endpoint is not set while deploying/updating/removing a stack, the program will try to guess it.
- Log messages contain a main message field and may contain several fields with context details, like stack name, endpoint name, warning implications, error fixing suggestions, etc.
- A Custom User-Agent header is sent on requests to the Portainer server to identify the client.
- Supported platforms and architectures linux 32/64 bit, darwin 32/64 bit, windows 32/64 bit, and arm7 32/64 bit.
- `completion` command to print Bash completion script.
- `config set` command to set configuration options.
- `config get` command to get configuration options.
- `config list|ls` command to print configuration options.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
- `container access` command to set access control for containers.
  - `--admins` flag to limit access to administrators.
  - `--private` flag to limit access to current user.
  - `--public` flag to give access to all users.
- `endpoint list|ls` command to print endpoints.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
- `endpoint group inspect` command to print endpoint group info.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
- `endpoint group list|ls` command to print endpoint groups.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
- `endpoint inspect` command to print endpoint info.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
- `help` command to print global help.
- `login` command to authenticate against Portainer server.
  - `--print` flag to print the retrieved auth token.
- `network access` command to set access control for networks.
  - `--admins` flag to limit access to administrators.
  - `--private` flag to limit access to current user.
  - `--public` flag to give access to all users.
- `secret access` command to set access control for secrets.
  - `--admins` flag to limit access to administrators.
  - `--private` flag to limit access to current user.
  - `--public` flag to give access to all users.
- `service access` command to set access control for services.
  - `--admins` flag to limit access to administrators.
  - `--private` flag to limit access to current user.
  - `--public` flag to give access to all users.
- `stack access` command to set access control for stacks.
  - `--admins` flag to limit access to administrators.
  - `--private` flag to limit access to current user.
  - `--public` flag to give access to all users.
- `stack deploy|up|create` command to deploy/update a stack.
  - `--endpoint` flag to set the endpoint to use.
  - `-e, --env-file` flag to set the file with environment variables to use with the stack.
  - `-r, --prune` flag to remove services that are no longer referenced.
  - `--replace-env` flag to replace environment variables instead of merging them while updating a stack.
  - `-c, --stack-file` flag to set the file with the YAML definition of the stack.
- `stack inspect` command to print stack info.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
  - `--endpoint` flag to filter stack by endpoint name.
- `stack list|ls` command to print stacks.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
  - `--endpoint` flag to filter stacks by endpoint name.
- `stack remove|rm|down` command to remove a stack.
  - `--endpoint` flag to set the endpoint to use.
  - `--strict` flag to fail if the stack does not exist.
- `status` command to show Portainer server status.
  - `--format` flag to select output format from "table", "json" or a custom Go template. Defaults to "table".
- `volume access` command to set access control for volumes.
  - `--admins` flag to limit access to administrators.
  - `--private` flag to limit access to current user.
  - `--public` flag to give access to all users.
- `-h, --help` flags on each command to print its help.
- `-A, --auth-token` global flag to set Portainer auth token.
- `--config` global flag to set the path to a configuration file. Supported file formats are JSON, TOML, YAML, HCL, envfile and Java properties config files. Defaults to "$HOME/.psu.yaml".
- `-h, --help` global flag to print global help.
- `--log-format` global flag to set log format from "text" and "json". Defaults to "text".
- `-v, --log-level` global flag to set log level from "panic", "faltal", "error", "warning", "info", "debug" and "trace". Defaults to "info".
- `--password` long name for `-p` global flag.
- `-t, --timeout` global flag to set a timeout for requests execution.
- `--url` long name for `-l` global flag.
- `--user` long name for `-u` global flag.
- `--version` global flag to print the program version. It includes the version number (major.minor.patch), the commit hash it was built from, the platform and architecture it was compiled for, and the build date.

### Removed
- `-a` global flag to select an action to execute. The sytax is now `COMMAND ARG --FLAG`, with a command for each action.
- Verbose and debug mode, which used to be enabled through `-v` and `-d` global flags respectively. The Debug and Trace log levels are the new equivalents.
- Requirement off `bash`, `jq` and `httpie` to run the program (they have also been removed from the Docker image). The new executable binary is self-sufficient.

### Changed
- Single executable binary instead of a bash script. The program has been rewritten in Go language.
- Supported Portainer API raised to `1.22.0`. **Previous versions may still work.**
- Endpoints are refered to by their name instead of their id.
- Base Docker image fixed to `alpine:3.10`.
- `-s` global flag to disable SSL certificate validation renamed to `-i, --insecure`. It does not receive a value anymore, it is a boolean flag.
- `-n` global flag to set the stack name moved as a parameter to the `stack deploy`, `stack list` and `stack remove` commands.
- `-c` global flag to set the file with the YAML definition of the stack renamed to `-c, --stack-file` and moved to the `stack deploy` command.
- `-g` global flag to set the file with the environment variables to use with the stack renamed to `-e, --env-file` and moved to the `stack deploy` command.
- `-r` global flag to remove no longer referenced services renamed to `-p, --prune` and moved to the `stack deploy` command.
- `-e` global flag renamed to `--endpoint` and moved to the `stack deploy`, `stack list` and `stack remove` commands. It now expects an endpoint name instead of its id.
- `-t` global flag renamed to `--strict` and moved to the `stack remove` command. It does not receive a value anymore, it is a boolean flag.
- All supported environment variables prefixed with "PSU_" and renamed to match command and flag names.

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
