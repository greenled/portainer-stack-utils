#!/usr/bin/env bash
set -e
[[ "$TRACE" ]] && set -x

function registry_login() {
  if [[ -n "$CI_REGISTRY_USER" ]]; then
    echo "Logging to GitLab Container Registry with CI credentials..."
    docker login -u "$CI_REGISTRY_USER" -p "$CI_REGISTRY_PASSWORD" "$CI_REGISTRY"
    echo ""
  fi
}

function external_registry_login() {
  if [[ -n "$DOCKER_USER" ]]; then
    echo "Logging to External Registry..."
    docker login -u "$DOCKER_USER" -p "$DOCKER_PASSWORD" "$DOCKER_REGISTRY"
    echo ""
  fi
}

function setup_docker() {
  if ! docker info &>/dev/null; then
    if [ -z "$DOCKER_HOST" -a "$KUBERNETES_PORT" ]; then
      export DOCKER_HOST='tcp://localhost:2375'
    fi
  fi
}

function git_tag_on_success() {
  local git_tag="${1:-dev}"
  local target_branch="${2:-master}"

  if (
    [ "$CI_COMMIT_REF_NAME" == "$target_branch" ] &&
    [ -n "$GITLAB_API_TOKEN" ] &&
    [ -z "$GIT_RESET_TAG" ]
  ); then
    # (re)write Protected Tag
    # TODO: rewrite these 'wget' commands with 'curl' commands
    wget -Y off -O response.txt --header='Accept-Charset: UTF-8' --header "PRIVATE-TOKEN: $GITLAB_API_TOKEN" --post-data '_method=delete' $CI_API_V4_URL/projects/$CI_PROJECT_ID/protected_tags/$git_tag || true
    wget -Y off -O response.txt --header='Accept-Charset: UTF-8' --header "PRIVATE-TOKEN: $GITLAB_API_TOKEN" --post-data '_method=delete' $CI_API_V4_URL/projects/$CI_PROJECT_ID/repository/tags/$git_tag || true
    wget -Y off -O response.txt --header='Accept-Charset: UTF-8' --header "PRIVATE-TOKEN: $GITLAB_API_TOKEN" --post-data "tag_name=$git_tag&ref=$CI_COMMIT_SHA" $CI_API_V4_URL/projects/$CI_PROJECT_ID/repository/tags
    wget -Y off -O response.txt --header='Accept-Charset: UTF-8' --header "PRIVATE-TOKEN: $GITLAB_API_TOKEN" --post-data "name=$git_tag&create_access_level=0" $CI_API_V4_URL/projects/$CI_PROJECT_ID/protected_tags
  else
    echo WARNING: \$GITLAB_API_TOKEN variable is missing
  fi
}

function registry_tag_on_success() {
  local current_registry_tag="${1:-$CI_COMMIT_SHA}"
  local target_registry_tag="${2:-dev}"
  local target_branch="${3:-master}"
  local current_registry_image="${4:-$CI_REGISTRY_IMAGE/builds}"
  local target_registry_image="${5:-$CI_REGISTRY_IMAGE}"
  local target_external_registry_image="${6:-$DOCKER_REGISTRY_IMAGE}"

  if [ "$CI_COMMIT_REF_NAME" == "$target_branch" ]; then
    docker pull "$current_registry_image:$current_registry_tag"
    docker tag "$current_registry_image:$current_registry_tag" "$target_registry_image:$target_registry_tag"
    docker push "$target_registry_image:$target_registry_tag"
    if [ -n "$target_external_registry_image" ]; then
      docker tag "$current_registry_image:$current_registry_tag" "$target_external_registry_image:$target_registry_tag"
      docker push "$target_external_registry_image:$target_registry_tag"
    fi
  fi
}

# Reset the git repository to the target tag
#
# First argument pass to this function or the `GIT_RESET_TAG` CI variable
# must be set
#   git_reset_from_tag dev
# or:
#   GIT_RESET_TAG=dev
#   git_reset_from_tag
function git_reset_from_tag() {
  local git_target_tag="${1:-$GIT_RESET_TAG}"

  if (
    [ "$CI_PIPELINE_SOURCE" == "schedule" ] &&
    [ -n "$git_target_tag" ] && [ "$GIT_STRATEGY" != "none" ] &&
    [ -z "$CI_COMMIT_TAG" ]
  ); then
    # Get specific tag
    git reset --hard $git_target_tag
    export CI_COMMIT_SHA=$(git rev-parse HEAD)
    export CI_COMMIT_SHORT_SHA=$(git rev-parse --short HEAD)
  else
    echo NOTICE: Not a Scheduling Pipeline, skip the git tag reset stuff... # debug
  fi
}

# Get latest stable semantic versioning git tag
# from a specific git branch
#
# First argument pass to this function or the `CI_COMMIT_REF_NAME` CI variable
# must be set
#   get_git_last_stable_tag 1-0-stable
#   -> "v1.0.3"
# or:
#   CI_COMMIT_REF_NAME=1-0-stable
#   get_git_last_stable_tag
#   -> "v1.0.3"
#
# see: https://semver.org
function get_git_last_stable_tag() {
  local target_branch="${1:-$CI_COMMIT_REF_NAME}"

  git fetch origin $target_branch
  git checkout -f -q $target_branch
  echo "$(git tag --merged $target_branch | grep -w '^v[0-9]\+\.[0-9]\+\.[0-9]\+$' | sort -r -V | head -n 1)"
}

# Useful for updating Docker images, on release/stable branches, but not the psu code
# See: https://docs.gitlab.com/ce/workflow/gitlab_flow.html#release-branches-with-gitlab-flow
# You can create a scheduled pipeline with a targeted git branch ("master", "1-0-stable", ...)
# and the CI variables below:
#   "GIT_RESET_LAST_STABLE_TAG=true"
#   "DOCKER_CACHE_DISABLED=true"
#   "TEST_DISABLED=" # no value to unset this variable
# See: https://gitlab.com/help/user/project/pipelines/schedules
function git_reset_from_last_stable_tag() {
  if [ "$GIT_RESET_LAST_STABLE_TAG" == "true" ]; then
    git_last_stable_tag="$(get_git_last_stable_tag)"
    if [ -n "$git_last_stable_tag" ]; then
      export CI_COMMIT_REF_PROTECTED="true"
      export CI_COMMIT_TAG="$git_last_stable_tag"
      export GIT_RESET_TAG="$git_last_stable_tag"
      git_reset_from_tag
    else
      echo WARNING: Last stable git tag not found
    fi
  fi
}
