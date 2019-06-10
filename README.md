# Portainer Stack Utils

[![Docker Automated build](https://img.shields.io/docker/automated/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Docker Pulls](https://img.shields.io/docker/pulls/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Microbadger](https://images.microbadger.com/badges/image/greenled/portainer-stack-utils.svg)](http://microbadger.com/images/greenled/portainer-stack-utils "Image size")

Bash script to deploy/update/undeploy stacks in a [Portainer](https://portainer.io/) instance from a [docker-compose](https://docs.docker.com/compose) [yaml file](https://docs.docker.com/compose/compose-file). Based on previous work by [@vladbabii](https://github.com/vladbabii) on [docker-how-to/portainer-bash-scripts](https://github.com/docker-how-to/portainer-bash-scripts).

## Supported Portainer API

Script was created for the latest Portainer API, which at the time of writing is [1.19.2](https://app.swaggerhub.com/apis/deviantony/Portainer/1.19.2).

## How to install

Just clone the repo and use the script

```bash
git clone https://github.com/greenled/portainer-stack-utils.git
cd portainer-stack-utils
./psu ...
```

### Requirements

You will need these dependecies installed:

- [bash](https://www.gnu.org/software/bash/)
- [httpie](https://httpie.org/)
- [jq](https://stedolan.github.io/jq/)

For Debian and similar apt-powered systems: `apt install bash httpie jq`.

## How to use

The provided `psu` script allows to deploy/update/undeploy Portainer stacks. Settings can be passed through envvars and/or flags. Both envvars and flags can be mixed but flags will always overwrite envvar values. When deploying a stack, if it doesn't exist a new one is created, otherwise it's updated (unless strict mode is active).

### With envvars

This is particularly useful for CI/CD pipelines.

- `ACTION` ("deploy" or "undeploy", required): Whether to deploy or undeploy the stack
- `PORTAINER_USER` (string, required): Username
- `PORTAINER_PASSWORD` (string, required): Password
- `PORTAINER_URL` (string, required): URL to Portainer
- `PORTAINER_STACK_NAME` (string, required): Stack name
- `DOCKER_COMPOSE_FILE` (string, required if action=deploy): Path to doker-compose file
- `ENVIRONMENT_VARIABLES_FILE` (string, optional, only used when action=deploy or action=update): Path to file with environment variables to be used by the stack. See [stack environment variables](#stack-environment-variables) below.
- `PORTAINER_PRUNE` ("true" or "false", optional): Whether to prune unused containers or not. Defaults to `"false"`.
- `PORTAINER_ENDPOINT` (int, optional): Which endpoint to use. Defaults to `1`.
- `HTTPIE_VERIFY_SSL` ("yes" or "no", optional): Whether to verify SSL certificate or not. Defaults to `"yes"`.
- `VERBOSE_MODE` ("true" or "false", optional): Whether to activate verbose output mode or not. Defaults to `"false"`. See [verbose mode](#verbose-mode) below.
- `DEBUG_MODE` ("true" or "false", optional): Whether to activate debug output mode or not. Defaults to `"false"`. See [debug mode](#debug-mode) below.
- `STRICT_MODE` ("true" or "false", optional): Whether to activate strict mode or not. Defaults to `"false"`. See [strict mode](#strict-mode) below.

#### Examples

```bash
export ACTION="deploy"
export PORTAINER_USER="admin"
export PORTAINER_PASSWORD="password"
export PORTAINER_URL="http://portainer.local"
export PORTAINER_STACK_NAME="mystack"
export DOCKER_COMPOSE_FILE="/path/to/docker-compose.yml"
export ENVIRONMENT_VARIABLES_FILE="/path/to/env_vars_file"

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
- `-g` (string, optional, only used when action=deploy or action=update): Path to file with environment variables to be used by the stack. See [stack environment variables](#stack-environment-variables) below.
- `-r` ("true" or "false", optional): Whether to prune unused containers or not. Defaults to `"false"`.
- `-e` (int, optional): Which endpoint to use. Defaults to `1`.
- `-s` ("yes" or "no", optional): Whether to verify SSL certificate or not. Defaults to `"yes"`.
- `-v` ("true" or "false", optional): Whether to activate verbose output mode or not. Defaults to `"false"`. See [verbose mode](#verbose-mode) below.
- `-d` ("true" or "false", optional): Whether to activate debug output mode or not. Defaults to `"false"`. See [debug mode](#debug-mode) below.
- `-t` ("true" or "false", optional): Whether to activate strict mode or not. Defaults to `"false"`. See [strict mode](#strict-mode) below.

#### Examples

```bash
./psu -a deploy -u admin -p password -l http://portainer.local -n mystack -c /path/to/docker-compose.yml -g /path/to/env_vars_file
```

```bash
./psu -a undeploy -u admin -p password -l http://portainer.local -n mystack
```

### Stack environment variables

There can be set environment variables for each stack, be it a new deployment or an update. For example:

```bash
touch .env
echo "MYSQL_ROOT_PASSWORD=agoodpassword" >> .env
echo "ALLOWED_HOSTS=*" >> .env
./psu -a deploy -u admin -p password -l http://portainer.local -n django-stack -c /path/to/docker-compose.yml -g env_vars
```

Stack environment variables can be enabled through [ENVIRONMENT_VARIABLES_FILE envvar](#with-envvars) or [-g flag](#with-flags).

### Verbose mode

In verbose mode the script prints execution steps.

```text
Getting auth token...
Getting stack mystack...
Stack mystack not found.
Getting Docker info...
Getting swarm cluster (if any)...
Swarm cluster found.
Preparing stack JSON...
Creating stack mystack...
```

Verbose mode can be enabled through [VERBOSE_MODE envvar](#with-envvars) or [-v flag](#with-flags).

### Debug mode

In debug mode the script prints as much information as possible to help diagnosing a malfunction.

**WARNING**: Debug mode will print configuration values (with Portainer credentials) and Portainer API responses (with sensitive information like authentication token and stacks environment variables). Avoid using debug mode in CI/CD pipelines, as pipeline logs are usually recorded.

Debug mode can be enabled through [DEBUG_MODE envvar](#with-envvars) or [-d flag](#with-flags).

### Strict mode

In strict mode the script never updates an existent stack nor removes an unexistent one, and instead exits with an error.

Strict mode can be enabled through [STRICT_MODE envvar](#with-envvars) or [-t flag](#with-flags).

## License

Source code contained by this project is licensed under the [GNU General Public License version 3](https://www.gnu.org/licenses/gpl-3.0.en.html). See [LICENSE](LICENSE) file for reference.
