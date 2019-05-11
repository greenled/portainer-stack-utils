FROM alpine

ENV LANG="en_US.UTF-8" \
  LC_ALL="C.UTF-8" \
  LANGUAGE="en_US.UTF-8" \
  TERM="xterm" \
  ACTION="" \
  PORTAINER_USER="" \
  PORTAINER_PASSWORD="" \
  PORTAINER_URL="" \
  PORTAINER_STACK_NAME="" \
  DOCKER_COMPOSE_FILE="" \
  ENVIRONMENT_VARIABLES_FILE="" \
  PORTAINER_PRUNE="false" \
  PORTAINER_ENDPOINT="1" \
  HTTPIE_VERIFY_SSL="yes" \
  VERBOSE_MODE="false" \
  DEBUG_MODE="false" \
  STRICT_MODE="false"

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
