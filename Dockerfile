FROM alpine

ENV LANG="en_US.UTF-8" \
  LC_ALL="C.UTF-8" \
  LANGUAGE="en_US.UTF-8" \
  TERM="xterm" \
  ACTION="" \
  PORTAINER_USER="root" \
  PORTAINER_PASSWORD="password" \
  PORTAINER_URL="http://example.com:9000" \
  PORTAINER_STACK_NAME="" \
  DOCKER_COMPOSE_FILE="" \
  PORTAINER_PRUNE="false" \
  PORTAINER_ENDPOINT="1"
  HTTPIE_VERIFY_SSL="yes" \
  VERBOSE_MODE="false" \
  DEBUG_MODE="false"

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
