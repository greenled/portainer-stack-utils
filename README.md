# Portainer Stack Utils

[![CircleCI](https://circleci.com/gh/greenled/portainer-stack-utils.svg?style=svg)](https://circleci.com/gh/greenled/portainer-stack-utils)
[![Docker Automated build](https://img.shields.io/docker/automated/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Docker Pulls](https://img.shields.io/docker/pulls/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Microbadger](https://images.microbadger.com/badges/image/greenled/portainer-stack-utils.svg)](http://microbadger.com/images/greenled/portainer-stack-utils "Image size")
[![Go Report Card](https://goreportcard.com/badge/github.com/greenled/portainer-stack-utils)](https://goreportcard.com/report/github.com/greenled/portainer-stack-utils)

## Overview

Portainer Stack Utils is a CLI client for [Portainer](https://portainer.io/) written in Go.

## Supported Portainer API

This application was created for the latest Portainer API, which at the time of writing is [1.22.0](https://app.swaggerhub.com/apis/deviantony/Portainer/1.22.0).

## How to install

Download the binaries for your platform from [the releases page](https://github.com/greenled/portainer-stack-utils/releases). The binaries have no external dependencies.

You can also install the source code with `go` and build the binaries yourself.

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

Each flag can be set inline (i.e. `--debug`), through an environment variable (i.e. `PSU_DEBUG=true`) and through a configuration file ([see below](#with-configuration-file)). All three methods can be combined, but if a flag is set more than once the order of precedence is:

1. Inline flag
2. Environment variable
3. Configuration file

#### With inline flags

Each command has it's own flags. Run `psu [COMMAND [SUBCOMMAND]] --help` to see each command's flag set.

```bash
psu --help
psu stack --help
psu stack deploy --help
```

#### With environment variables

This is particularly useful for CI/CD pipelines.

Environment variables can be bound to flags following the `PSU_[COMMAND_[SUBCOMMAND_]]FLAG` naming pattern:

| Command and subcommand | Flag | Environment variable | Comment |
| :--------------------- | :--- | :------------------- | :------ |
|  | --verbose | PSU_VERBOSE=true | All environment variables are prefixed with "PSU_" |
| stack list | --quiet | PSU_STACK_LIST_QUIET=true | Commands and subcommands are uppercased and joined with "_" |
| stack deploy | --env-file .env | PSU_STACK_DEPLOY_ENV_FILE=.env | Characters "-" in flag name are replaced with "_" |

#### With configuration file

Flags can be bound to a configuration file. Use the `--config` flag to specify a configuration file to load flags from. By default the file `$HOME/.psu.yaml` is used if present.

#### Using Yaml

If you use a Yaml configuration file:

```text
[command:
  [subcommand:]]
    flag: value
```

```yaml
verbose: true
url: http://localhost:10000
insecure: true
stack:
  deploy:
    stack-file: docker-compose.yml
    env-file: .env
  list:
    quiet: true
```

This is valid too:

```text
[command.[subcommand.]]flag: value
```

```yaml
verbose: true
url: http://localhost:10000
insecure: true
stack.deploy.stack-file: docker-compose.yml
stack.deploy.env-file: .env
stack.list.quiet: true
```

#### Using Json

If you use a Json configuration file:

```text
{
  ["command": {
    ["subcommand": {]]
      "flag": value
    [}]
  [}]
{
```

```json
{
  "verbose": true,
  "url": "http://localhost:10000",
  "insecure": true,
  "stack": {
    "deploy": {
      "stack-file": "docker-compose.yml",
      "env-file": ".env"
    },
    "list": {
      "quiet": true
    }
  }
}
```

This is valid too:

```text
{
  "[command.[subcommand.]]flag": value
}
```

```json
{
"verbose": true,
"url": "http://localhost:10000",
"insecure": true,
"stack.deploy.stack-file": "docker-compose.yml",
"stack.deploy.env-file": ".env",
"stack.list.quiet": true
}

```

### Stack environment variables

You will usually want to set some environment variables in your stacks. You can do so with the `--env-file` flag:

```bash
touch .env
echo "MYSQL_ROOT_PASSWORD=agoodpassword" >> .env
echo "ALLOWED_HOSTS=*" >> .env
psu stack deploy django-stack -c /path/to/docker-compose.yml -e .env
```

As every flag, this one can also be used with the `PSU_STACK_DEPLOY_ENV_FILE` [environment variable](#with-environment-variables) and the `psu.stack.deploy.env-file` [configuration key](#with-configuration-file).

### Verbose mode

In verbose mode the script prints execution steps.

```text
2019/07/20 19:15:45 [Using config file: /home/johndoe/.psu.yaml]
2019/07/20 19:15:45 [Getting stack mystack...]
2019/07/20 19:15:45 [Getting auth token...]
2019/07/20 19:15:45 [Stack mystack not found. Deploying...]
2019/07/20 19:15:45 [Swarm cluster found with id qwe123rty456uio789asd123f]
```

Verbose mode can be enabled through the `PSU_VERBOSE` [environment variable](#with-environment-variables) and the `verbose` [configuration key](#with-configuration-file).

### Debug mode

In debug mode the script prints as much information as possible to help diagnosing a malfunction.

**WARNING**: Debug mode will print configuration values (with Portainer credentials) and Portainer API responses (with sensitive information like authentication token and stacks environment variables). Avoid using debug mode in CI/CD pipelines, as pipeline logs are usually recorded.

Debug mode can be enabled through the `PSU_DEBUG` [environment variable](#with-environment-variables) and the `debug` [configuration key](#with-configuration-file).

## Contributing

So, you want to contribute to the project:

- Fork it
- Download your fork to your PC (git clone https://github.com/your_username/portainer-stack-utils && cd portainer-stack-utils)
- Create your feature branch (git checkout -b my-new-feature)
- Make changes and add them (git add .)
- Commit your changes (git commit -m 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create a new pull request

If you are submitting a complex feature, create a small design proposal on the [issue tracker](https://github.com/greenled/portainer-stack-utils/issues) before you start.

## License

Source code contained by this project is licensed under the [GNU General Public License version 3](https://www.gnu.org/licenses/gpl-3.0.en.html). See [LICENSE](LICENSE) file for reference.
