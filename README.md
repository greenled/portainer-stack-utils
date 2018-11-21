# Portainer Stack Utils

[![Docker Automated build](https://img.shields.io/docker/automated/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Docker Pulls](https://img.shields.io/docker/pulls/greenled/portainer-stack-utils.svg)](https://hub.docker.com/r/greenled/portainer-stack-utils/)
[![Microbadger](https://images.microbadger.com/badges/image/greenled/portainer-stack-utils.svg)](http://microbadger.com/images/greenled/portainer-stack-utils "Image size")

Bash scripts to deploy/undeploy stacks in a [Portainer](https://portainer.io/) instance from a [docker-compose](https://docs.docker.com/compose) [yaml file](https://docs.docker.com/compose/compose-file). Based on previous work by [@vladbabii](https://github.com/vladbabii) on [docker-how-to/portainer-bash-scripts](https://github.com/docker-how-to/portainer-bash-scripts).

## Requirements

- [bash](https://www.gnu.org/software/bash/)
- [httpie](https://httpie.org/)
- [jq](https://stedolan.github.io/jq/)

## How to use

Two scripts are included: `deploy.sh` and `undeploy.sh`. Both scripts use the following environment variables to connect to the portainer instance:

- `PORTAINER_USER` (string): Username
- `PORTAINER_PASSWORD` (string): Password
- `PORTAINER_URL` (string): URL to Portainer
- `PORTAINER_PRUNE` ("true" or "false"): Whether to prune unused containers or not. Defaults to `"false"`.
- `PORTAINER_ENDPOINT` (int): Which endpoint to use. Defaults to `1`.
- `HTTPIE_VERIFY_SSL` ("yes" or "no"): Whether to verify SSL certificate or not. Defaults to `"yes"`.

### deploy.sh

This script deploys a stack. The stack is created if it does not exist, otherwise it is updated. You must pass the stack name and the path to the docker-compose file as arguments:

```bash
./deploy.sh mystack docker-compose.yml
```

### undeploy.sh

This script removes a stack. You must pass the stack name as argument:

```bash
./undeploy.sh mystack
```
