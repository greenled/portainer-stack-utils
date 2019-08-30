# Portainer Stack Utils

[![CircleCI](https://circleci.com/gh/greenled/portainer-stack-utils.svg?style=svg)](https://circleci.com/gh/greenled/portainer-stack-utils)
[![Docker Automated build](https://img.shields.io/docker/automated/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Docker Pulls](https://img.shields.io/docker/pulls/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Microbadger](https://images.microbadger.com/badges/image/greenled/portainer-stack-utils.svg)](http://microbadger.com/images/greenled/portainer-stack-utils "Image size")
[![Go Report Card](https://goreportcard.com/badge/github.com/greenled/portainer-stack-utils)](https://goreportcard.com/report/github.com/greenled/portainer-stack-utils)

## Table of contents

- [Overview](#overview)
- [Supported Portainer API](#supported-portainer-api)
- [How to install](#how-to-install)
- [How to use](#how-to-use)
  - [Configuration](#configuration)
    - [Inline flags](#inline-flags)
    - [Environment variables](#environment-variables)
    - [Configuration files](#configuration-files)
      - [YAML configuration file](#yaml-configuration-file)
      - [JSON configuration file](#json-configuration-file)
  - [Environment variables for deployed stacks](#environment-variables-for-deployed-stacks)
  - [Endpoint's Docker API proxy](#endpoints-docker-api-proxy)
  - [Log level](#log-level)
  - [Exit statuses](#exit-statuses)
- [Contributing](#contributing)
- [License](#license)

## Overview

Portainer Stack Utils is a CLI client for [Portainer](https://portainer.io/) written in Go.

**Attention:** The `master` branch contains the next major version, still unstable and under heavy development. A more stable (and also older) version is available as a Bash script in [release 0.1.1](https://github.com/greenled/portainer-stack-utils/releases/0.1.1), and also as a [Docker image](https://hub.docker.com/r/greenled/portainer-stack-utils). There is ongoing work in `1-0-next` branch to enhace that Bash version.

## Supported Portainer API

This application was created for the latest Portainer API, which at the time of writing is [1.22.0](https://app.swaggerhub.com/apis/deviantony/Portainer/1.22.0).

## How to install

Download the binaries for your platform and architecture from [the releases page](https://github.com/greenled/portainer-stack-utils/releases).

## How to use

The application is built on a structure of commands, arguments and flags.
                   
**Commands** represent actions, **Args** are things and **Flags** are modifiers for those actions:

```text
APPNAME COMMAND ARG --FLAG
```

Here are some examples:

```bash
psu help
psu status --help
psu stack ls --endpoint primary --format "{{ .Name }}"
psu stack deploy mystack --stack-file docker-compose.yml -e .env --log-level debug
psu stack rm mystack
```

Commands can have subcommands, like `stack ls` and `stack deploy` in the previous example. They can also have aliases (i.e. `create` and `up` are aliases of `deploy`).

Some flags are global, which means they affect every command (i.e. `--log-level`), while others are local, which mean they only affect the command they belong to (i.e. `--stack-file` flag from `deploy` command). Also, some flags have a short version (i.e `--insecure`, `-i`).

### Configuration

The program can be configured through [inline flags](#inline-flags) (i.e. `--user`), [environment variables](#environment-variables) (i.e. `PSU_USER=admin`) and/or [configuration files](#configuration-files), which translate into multi-level configuration keys in the form `x[.y[.z[...]]]`. Run `psu config ls` to see all available configuration options.

All three methods can be combined. If a configuration key is set in several places the order of precedence is:

1. Inline flags
2. Environment variables
3. Configuration file
4. Default values

#### Inline flags

Configuration can be set through inline flags. Valid combinations of commands and flags directly map to configuration keys:

| Configuration key | Command | Flag |
| :---------------- | :------ | :--- |
| user | psu | --user |
| stack.list.format | psu stack list | --format |
| stack.deploy.env-file | stack deploy | --env-file |

Run `psu help COMMAND` to see all available flags for a given command.

#### Environment variables

Configuration can be set through environment variables. Supported environment variables follow the `PSU_[COMMAND_[SUBCOMMAND_]]FLAG` naming pattern:

| Configuration key | Environment variable |
| :---------------- | :------------------- |
| user | PSU_USER |
| stack.list.format | PSU_STACK_LIST_FORMAT |
| stack.deploy.env-file | PSU_STACK_DEPLOY_ENV_FILE |

*Note that all supported environment variables are prefixed with "PSU_" to avoid name clashing. Characters "-" and "." in configuration key names are replaced with "_" in environment variable names.*

#### Configuration files

Configuration can be set through a configuration file. Supported file formats are [JSON](#json-configuration-file), TOML, [YAML](#yaml-configuration-file), HCL, envfile and Java properties config files. Use the `--config` global flag to specify a configuration file. File `$HOME/.psu.yaml` is used by default if present.

##### YAML configuration file

A Yaml configuration file should look like this:

```yaml
log-level: debug
user: admin
insecure: true
stack.list.format: table
stack:
  deploy.env-file: .env
  deploy:
    stack-file: docker-compose.yml
```

*Note that flat and nested keys are both valid.*

##### JSON configuration file

A JSON configuration file should look like this:

```json
{
  "log-level": "debug",
  "user": "admin",
  "insecure": true,
  "stack.list.format": "table",
  "stack": {
    "deploy.env-file": ".env",
    "deploy": {
      "stack-file": "docker-compose.yml"
    }
  }
}
```

*Note that flat and nested keys are both valid.*

### Environment variables for deployed stacks

You will often want to set environment variables in your deployed stacks. You can do so through the `stack.deploy.env-file` [configuration key](#configuration). :

```bash
touch .env
echo "MYSQL_ROOT_PASSWORD=agoodpassword" >> .env
echo "ALLOWED_HOSTS=*" >> .env

# Using --env-file flag
psu stack deploy django-stack -c /path/to/docker-compose.yml -e .env

# Using PSU_STACK_DEPLOY_ENV_FILE environment variable
PSU_STACK_DEPLOY_ENV_FILE=.env
psu stack deploy django-stack -c /path/to/docker-compose.yml

# Using a config file
echo "stack.deploy.env-file: .env" > .config.yml
psu stack deploy django-stack -c /path/to/docker-compose.yml --config .config.yml
```

### Endpoint's Docker API proxy

If you want finer-grained control over an endpoint's Docker daemon you can expose it through a proxy and configure a local Docker client to use it.

First, expose the endpoint's Docker API:

```bash
psu proxy --endpoint primary --address 127.0.0.1:2375
```

Then (in a different shell), configure a local Docker client to use the exposed API:

```bash
export DOCKER_HOST=tcp://127.0.0.1:2375
```

Now you can run `docker ...` commands in the `primary` endpoint as in a local Docker installation, **with the added benefit of using Portainer's RBAC**.

*Note that creating stacks through `docker stack ...` instead of `psu stack ...` will give you *limited* control over them, as they are created outside of Portainer.*

### Log level

You can control how much noise you want the program to do by setting the log level. There are seven log levels:

- *panic*: Unexpected errors that stop program execution.
- *fatal*: Expected errors that stop program execution.
- *error*: Errors that should definitely be noted but don't stop the program execution.
- *warning*: Non-critical events that deserve eyes.
- *info*: General events about what's going on inside the program. This is the default level.
- *debug*: Very verbose logging. Usually only enabled when debugging.
- *trace*: Finer-grained logging than the *debug* level.

**WARNING**: **trace** level will print sensitive information, like Portainer API requests and responses (with authentication tokens, stacks environment variables, and so on). Avoid using **trace** level in CI/CD environments, as those logs are usually recorded.

This is an example with *debug* level:

```bash
psu stack deploy asd --endpoint primary --log-level debug
```

The output would look like:

```text
DEBU[0000] Getting endpoint's Docker info     endpoint=primary
DEBU[0000] Getting stack                      endpoint=primary stack=asd
DEBU[0000] Stack not found                    stack=asd
INFO[0000] Creating stack                     endpoint=primary stack=asd
INFO[0000] Stack created                      endpoint=primary id=89 stack=asd
```

### Exit statuses

- *0*: Program executed normally.
- *1*: An expected error stopped program execution.
- *2*: An unexpected error stopped program execution.

## Contributing

Contributing guidelines can be found in [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Source code contained by this project is licensed under the [GNU General Public License version 3](https://www.gnu.org/licenses/gpl-3.0.en.html). See [LICENSE](LICENSE) file for reference.
