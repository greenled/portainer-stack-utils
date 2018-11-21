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

if [ -z ${2+x} ]; then
  echo "Error: Parameter #2 missing (path to yaml)"
  exit 1
fi
STACK_YAML_PATH="$2"

STACK_YAML_CONTENT=$(cat "$STACK_YAML_PATH")

# Escape carriage returns
STACK_YAML_CONTENT="${STACK_YAML_CONTENT//$'\r'/''}"

# Escape double quotes
STACK_YAML_CONTENT="${STACK_YAML_CONTENT//$'"'/'\"'}"

# Escape newlines
STACK_YAML_CONTENT="${STACK_YAML_CONTENT//$'\n'/'\n'}"

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

  echo "Getting swarm cluster (if any)..."
  SWARM_ID=$(http \
    --ignore-stdin \
    --verify=$HTTPIE_VERIFY_SSL \
    "$PORTAINER_URL/api/endpoints/$PORTAINER_ENDPOINT/docker/info" \
    "Authorization: Bearer $AUTH_TOKEN" \
    | jq -r ".Swarm.Cluster.ID // empty")
  
  echo "Creating stack $STACK_NAME..."
  if [ -z "$SWARM_ID" ];then
    DATA_PREFIX="{\"Name\":\"$STACK_NAME\",\"StackFileContent\":\""
    DATA_SUFFIX="\"}"
    echo "$DATA_PREFIX$STACK_YAML_CONTENT$DATA_SUFFIX" > json.tmp

    CREATE=$(http \
      --ignore-stdin \
      --verify=$HTTPIE_VERIFY_SSL \
      --timeout=300 \
      "$PORTAINER_URL/api/stacks" \
      "Authorization: Bearer $AUTH_TOKEN" \
      type==2 \
      method==string \
      endpointId==$PORTAINER_ENDPOINT \
      @json.tmp)
  else
  	DATA_PREFIX="{\"Name\":\"$STACK_NAME\",\"SwarmID\":\"$SWARM_ID\",\"StackFileContent\":\""
    DATA_SUFFIX="\"}"
    echo "$DATA_PREFIX$STACK_YAML_CONTENT$DATA_SUFFIX" > json.tmp

    CREATE=$(http \
      --ignore-stdin \
      --verify=$HTTPIE_VERIFY_SSL \
      --timeout=300 \
      "$PORTAINER_URL/api/stacks" \
      "Authorization: Bearer $AUTH_TOKEN" \
      type==1 \
      method==string \
      endpointId==$PORTAINER_ENDPOINT \
      @json.tmp)
  fi

  rm json.tmp

  if [ -z ${CREATE+x} ]; then
    echo "Error: stack $STACK_NAME not created"
    exit 1
  fi
else
  echo "Result: Stack $STACK_NAME found."

  STACK_ID="$(echo "$STACK" | jq -j ".Id")"
  STACK_ENV_VARS="$(echo -n "$STACK"| jq ".Env" -jc)"
  DATA_PREFIX="{\"Id\":\"$STACK_ID\",\"StackFileContent\":\""
  DATA_SUFFIX="\",\"Env\":"$STACK_ENV_VARS",\"Prune\":$PORTAINER_PRUNE}"
  echo "$DATA_PREFIX$STACK_YAML_CONTENT$DATA_SUFFIX" > json.tmp
  
  echo "Updating stack $STACK_NAME..."
  UPDATE=$(http \
    --ignore-stdin \
    --verify=$HTTPIE_VERIFY_SSL \
    --timeout=300 \
    PUT "$PORTAINER_URL/api/stacks/$STACK_ID" \
    "Authorization: Bearer $AUTH_TOKEN" \
    endpointId==$PORTAINER_ENDPOINT \
    @json.tmp)
  
  rm json.tmp
  
  if [ -z ${UPDATE+x} ]; then
    echo "Error: stack $STACK_NAME not updated"
    exit 1
  fi
fi
echo "Done"
