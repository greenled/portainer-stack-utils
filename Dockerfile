FROM alpine:3.10

RUN apk add --no-cache \
  bash ca-certificates gettext httpie jq \
  py3-pip python3-dev libc-dev libffi-dev openssl-dev gcc make; \
  \
  pip3 --no-cache-dir install docker-compose; \
  \
  apk del python3-dev libc-dev libffi-dev openssl-dev gcc make; \
  rm -rf /tmp/src

ENV LANG="en_US.UTF-8" \
    LC_ALL="C.UTF-8" \
    LANGUAGE="en_US.UTF-8" \
    TERM="xterm" \
    ACTION="" \
    PORTAINER_USER="" \
    PORTAINER_PASSWORD="" \
    PORTAINER_AUTH_TOKEN="" \
    PORTAINER_URL="" \
    PORTAINER_STACK_NAME="" \
    PORTAINER_SERVICE_NAME="" \
    DOCKER_COMPOSE_FILE="" \
    DOCKER_COMPOSE_LINT="true" \
    ENVIRONMENT_VARIABLES_FILE="" \
    PORTAINER_ENDPOINT="1" \
    PORTAINER_PRUNE="false" \
    TIMEOUT=100 \
    AUTO_DETECT_JOB="true" \
    HTTPIE_VERIFY_SSL="yes" \
    VERBOSE_MODE="false" \
    DEBUG_MODE="false" \
    QUIET_MODE="false" \
    STRICT_MODE="false" \
    MASKED_VARIABLES="false"

COPY psu /usr/local/bin/

RUN chmod +x /usr/local/bin/*

ENTRYPOINT ["/usr/local/bin/psu"]
