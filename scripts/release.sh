#!/usr/bin/env bash
set -e
[[ "$TRACE" ]] && set -x

source "$(dirname "$0")/helpers.sh"

# First argument is the variant(s) of the Docker image to be tagged
# e.g. "core"
# or "core debian debian-core" for multiple variants
variant_input="$1"

# Split space (" ") character to new line "\n"
variants="$(echo "$variant_input" | tr ' ' '\n')"

if [ "$GIT_RESET_LAST_STABLE_TAG" != "true" ]; then
  git_tag_on_success "dev"
  registry_tag_on_success $CI_COMMIT_SHA "dev"
  for variant in $variants; do
    registry_tag_on_success "${variant}-${CI_COMMIT_SHA}" "dev-${variant}"
  done
fi

# Write current git tag to the Docker registry
# if it starts with the 'v' character, followed by a number, e.g 'v1.0.2'
# see: https://semver.org
# You should set a 'v*' protected tags, in your GitLab project
# see: https://gitlab.com/help/user/project/protected_tags#configuring-protected-tags
if (
  [[ "$CI_COMMIT_TAG" =~ ^v[0-9][0-9a-z.\-]*$ ]] &&
  [ "$CI_COMMIT_REF_PROTECTED" == "true" ]
); then
  # Remove the first letter ('v' character), for tagging Docker images
  target_registry_tag="${CI_COMMIT_TAG:1}"

  registry_tag_on_success $CI_COMMIT_SHA $target_registry_tag $CI_COMMIT_REF_NAME
  for variant in $variants; do
    registry_tag_on_success "${variant}-${CI_COMMIT_SHA}" "${target_registry_tag}-${variant}" $CI_COMMIT_REF_NAME
  done

  # If current git tag is a stable semantic version.
  # see: https://semver.org
  if [[ "$target_registry_tag" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    git fetch --all
    # e.g. "  origin/master"
    master_branch="$(git branch --no-color --remotes --list 'origin/master')"

    # e.g. "  origin/1-0-stable
    # >  origin/1-1-stable
    # >  origin/2-0-stable
    # >  origin/master"
    commit_in_branches="$(git branch --no-color --remotes --contains $CI_COMMIT_SHA | grep -w '  origin/master\|  origin/.*-stable')"

    # e.g. "  origin/2-0-stable
    # >  origin/1-1-stable
    # >  origin/1-0-stable"
    stable_branches="$(echo "$commit_in_branches" | grep -vw '  origin/master' | sort -rV)"

    # Extract the MAJOR and MAJOR.MINOR versions
    # Then tag and push them to Docker registry.
    # e.g. the "v1.3.7" git tag, will create/update "1" and "1.3" registry tags
    # For MAJOR tag on git stable branches, we use the latest stable branch
    # e.g. when there is 2 stable branches "1-0-stable" and "1-1-stable"
    # If there is a new minor git tag "v1.0.3", only on the "1-0-stable".
    # We don't have to update the registry tag  "1".
    # Because only git tags on the "1-1-stable" git branch can do it.
    major_registry_tag=$(echo $target_registry_tag | sed -E 's/^([0-9]+)\.[0-9]+\.[0-9]+$/\1/')
    minor_registry_tag=$(echo $target_registry_tag | sed -E 's/^[0-9]+\.([0-9]+)\.[0-9]+$/\1/')
    major_minor_registry_tag=$major_registry_tag.$minor_registry_tag
    target_stable_branch="${major_registry_tag}-${minor_registry_tag}-stable"

    # e.g. "1-1-stable"
    latest_major_branch="$(echo "$stable_branches"| grep -E "^  origin/$major_registry_tag-[0-9]+-stable$" | head -n 1 | sed -E 's/^  origin\/(.+)$/\1/')"

    # If the git tag is only contained in master branch
    # or if it's contained in the latest stable branch who has the same MAJOR version
    # e.g. when there is 2 stable branches "1-0-stable" and "1-1-stable".
    # If git tag is "v1.0.1", the MAJOR registry tag "1" is NOT written.
    # Because the lastet stable branch is "1-1-stable".
    #
    # If git tag is "v1.1.3", the MAJOR registry tag "1" is written.
    # Because the lastet stable branch is "1-1-stable".
    if (
      ([ -n "$master_branch" ] && [ "$master_branch" == "$commit_in_branches" ]) ||
      [ "$latest_major_branch" == "$target_stable_branch" ]
    ); then
      registry_tag_on_success $CI_COMMIT_SHA $major_registry_tag $CI_COMMIT_REF_NAME
    fi

    registry_tag_on_success $CI_COMMIT_SHA $major_minor_registry_tag $CI_COMMIT_REF_NAME
    for variant in $variants; do
      if (
        ([ -n "$master_branch" ] && [ "$master_branch" == "$commit_in_branches" ]) ||
        [ "$latest_major_branch" == "$target_stable_branch" ]
      ); then
        registry_tag_on_success "${variant}-${CI_COMMIT_SHA}" "${major_registry_tag}-${variant}" $CI_COMMIT_REF_NAME
      fi

      registry_tag_on_success "${variant}-${CI_COMMIT_SHA}" "${major_minor_registry_tag}-${variant}" $CI_COMMIT_REF_NAME
    done

    # The latest stable semantic versioning git tag on the "master" branch
    # is considered as the "latest" registry tag.
    # Check if the current git commit is only present in the "master" branch
    # and not in stable branches (e.g. "1-0-stable").
    # Non stable branches are skipped (e.g. "feature-better-performance").
    if [ -n "$master_branch" ] && [ "$master_branch" == "$commit_in_branches" ]; then
      registry_tag_on_success $CI_COMMIT_SHA "latest" $CI_COMMIT_REF_NAME
      for variant in $variants; do
        registry_tag_on_success "${variant}-${CI_COMMIT_SHA}" "${variant}" $CI_COMMIT_REF_NAME
      done
    fi
  fi
fi
