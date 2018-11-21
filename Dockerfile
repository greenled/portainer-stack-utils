FROM alpine

ENV LANG="en_US.UTF-8" \
  LC_ALL="C.UTF-8" \
  LANGUAGE="en_US.UTF-8" \
  TERM="xterm" \
  PORTAINER_USER="root" \
  PORTAINER_PASSWORD="password" \
  PORTAINER_URL="http://example.com:9000" \
  PORTAINER_PRUNE="false" \
  PORTAINER_ENDPOINT="1"

RUN apk --update add \
  bash \
  ca-certificates \
  httpie \
  jq \
  gettext \
  && \
  rm -rf /tmp/src && \
  rm -rf /var/cache/apk/*

RUN chmod +x *.sh

COPY *.sh /usr/local/bin
