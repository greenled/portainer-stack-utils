FROM alpine:3.10

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

RUN apk --update add \
  bash \
  ca-certificates \
  httpie \
  jq \
  gettext \
  && \
  rm -rf /tmp/src && \
  rm -rf /var/cache/apk/*

COPY psu /usr/local/bin/

RUN chmod +x /usr/local/bin/*

ENTRYPOINT ["/usr/local/bin/psu"]
