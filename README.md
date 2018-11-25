# Portainer Stack Utils

[![Docker Automated build](https://img.shields.io/docker/automated/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Docker Pulls](https://img.shields.io/docker/pulls/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Microbadger](https://images.microbadger.com/badges/image/greenled/portainer-stack-utils.svg)](http://microbadger.com/images/greenled/portainer-stack-utils "Image size")

Bash script to deploy/update/undeploy stacks in a [Portainer](https://portainer.io/) instance from a [docker-compose](https://docs.docker.com/compose) [yaml file](https://docs.docker.com/compose/compose-file). Based on previous work by [@vladbabii](https://github.com/vladbabii) on [docker-how-to/portainer-bash-scripts](https://github.com/docker-how-to/portainer-bash-scripts).

## Supported Portainer API

Script was created for the latest Portainer API, which at the time of writing is [1.9.2](https://app.swaggerhub.com/apis/deviantony/Portainer/1.19.2).

## Requirements

- [bash](https://www.gnu.org/software/bash/)
- [httpie](https://httpie.org/)
- [jq](https://stedolan.github.io/jq/)

## How to use

The provided `psu` script allows to deploy/update/undeploy Portainer stacks. Settings can be passed through envvars and/or flags. Both envvars and flags can be mixed but flags will always overwrite envvar values. When deploying a stack, if it doesn't exist a new one is created, otherwise it's updated.

### With envvars

This is particularly useful for CI/CD pipelines.

- `ACTION` ("deploy" or "undeploy", required): Whether to deploy or undeploy the stack
- `PORTAINER_USER` (string, required): Username
- `PORTAINER_PASSWORD` (string, required): Password
- `PORTAINER_URL` (string, required): URL to Portainer
- `PORTAINER_STACK_NAME` (string, required): Stack name
- `DOCKER_COMPOSE_FILE` (string, required if action=deploy): Path to doker-compose file
- `PORTAINER_PRUNE` ("true" or "false", optional): Whether to prune unused containers or not. Defaults to `"false"`.
- `PORTAINER_ENDPOINT` (int, optional): Which endpoint to use. Defaults to `1`.
- `HTTPIE_VERIFY_SSL` ("yes" or "no", optional): Whether to verify SSL certificate or not. Defaults to `"yes"`.
- `VERBOSE_MODE` ("true" or "false", optional): Whether to activate verbose output mode or not. Defaults to `"false"`.
- `DEBUG_MODE` ("true" or "false", optional): Whether to activate debug output mode or not. Defaults to `"false"`. See [debug mode warning](#debug-mode) below.

#### Examples

```bash
export ACTION="deploy"
export PORTAINER_USER="admin"
export PORTAINER_PASSWORD="password"
export PORTAINER_URL="http://portainer.local"
export PORTAINER_STACK_NAME="mystack"
export DOCKER_COMPOSE_FILE="/path/to/docker-compose.yml"

./psu
```

```bash
export ACTION="undeploy"
export PORTAINER_USER="admin"
export PORTAINER_PASSWORD="password"
export PORTAINER_URL="http://portainer.local"
export PORTAINER_STACK_NAME="mystack"

./psu
```

### With flags

This is more suitable for standalone script usage.

- `-a` ("deploy" or "undeploy", required): Whether to deploy or undeploy the stack
- `-u` (string, required): Username
- `-p` (string, required): Password
- `-l` (string, required): URL to Portainer
- `-n` (string, required): Stack name
- `-c` (string, required if action=deploy): Path to doker-compose file
- `-r` ("true" or "false", optional): Whether to prune unused containers or not. Defaults to `"false"`.
- `-e` (int, optional): Which endpoint to use. Defaults to `1`.
- `-s` ("yes" or "no", optional): Whether to verify SSL certificate or not. Defaults to `"yes"`.
- `-v` ("true" or "false", optional): Whether to activate verbose output mode or not. Defaults to `"false"`.
- `-d` ("true" or "false", optional): Whether to activate debug output mode or not. Defaults to `"false"`. See [debug mode warning](#debug-mode) below.

#### Examples

```bash
./psu -a deploy -u admin -p password -l http://portainer.local -n mystack -c /path/to/docker-compose.yml
```

```bash
./psu -a undeploy -u admin -p password -l http://portainer.local -n mystack
```

### Debug mode

**WARNING**: In debug mode the script prints as much information as possible, including configuration values (with Portainer credentials) and Portainer API responses (with sensitive information like authentication token and stacks environment variables). Avoid using debug mode in CI/CD pipelines, as pipeline logs are usually recorded.

Debug mode can be enabled through [DEBUG_MODE envvar](#with-envvars) or [-d flag](with-flags).

## License

Source code contained by this project is licensed under the [GNU General Public License version 3](https://www.gnu.org/licenses/gpl-3.0.en.html). See [LICENSE](LICENSE) file for reference.
