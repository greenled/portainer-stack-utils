#!/usr/bin/env bash

PORTAINER_USER=${PORTAINER_USER:-"user"}
PORTAINER_PASSWORD=${PORTAINER_PASSWORD:-"password"}
PORTAINER_URL=${PORTAINER_URL:-"https://portainer.example.com"}
PORTAINER_PRUNE=${PORTAINER_PRUNE:-"false"}
PORTAINER_ENDPOINT=${PORTAINER_ENDPOINT:-"1"}
HTTPIE_VERIFY_SSL=${HTTPIE_VERIFY_SSL:-"yes"}

if [ -z ${1+x} ]; then
  echo "Error: Parameter #1 missing (stack name)"
  exit 1
fi
STACK_NAME="$1"

echo "Getting auth token..."
AUTH_TOKEN=$(http \
  --ignore-stdin \
  --verify=$HTTPIE_VERIFY_SSL \
  $PORTAINER_URL/api/auth \
  username=$PORTAINER_USER \
  password=$PORTAINER_PASSWORD \
  | jq -r .jwt)

if [ -z "$AUTH_TOKEN" ]; then
  echo "Error: Authentication error."
  exit 1
fi
echo "Done"

echo "Getting stack $STACK_NAME..."
STACKS=$(http \
  --ignore-stdin \
  --verify=$HTTPIE_VERIFY_SSL \
  "$PORTAINER_URL/api/stacks" \
  "Authorization: Bearer $AUTH_TOKEN")

STACK=$(echo "$STACKS" \
  | jq --arg STACK_NAME "$STACK_NAME" -jc '.[] | select(.Name == $STACK_NAME)')

if [ -z "$STACK" ]; then
  echo "Result: Stack $STACK_NAME not found."
  exit 1
fi
echo "Result: Stack $STACK_NAME found."

STACK_ID="$(echo "$STACK" | jq -j ".Id")"

echo "Deleting stack $STACK_NAME..."
DELETE=$(http \
  --ignore-stdin \
  --verify=$HTTPIE_VERIFY_SSL \
  DELETE "$PORTAINER_URL/api/stacks/$STACK_ID" \
  "Authorization: Bearer $AUTH_TOKEN")

if [ -z ${DELETE+x} ]; then
  echo "Error: stack $STACK_NAME not deleted"
  exit 1
fi
echo "Done"
