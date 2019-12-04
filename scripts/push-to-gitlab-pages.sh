#!/usr/bin/env bash
# Usage:
# bash push-to-gitlab-pages.sh <target_branch_name> <path_to_push_in_the_branch>
# bash push-to-gitlab-pages.sh
# bash push-to-gitlab-pages.sh "gitlab-pages" "public"

set -e
[[ "$TRACE" ]] && set -x

if [ -z "$GITLAB_WRITE_REPO_USER" ] || [ -z "$GITLAB_WRITE_REPO_TOKEN" ]; then
  echo "ERROR: \$GITLAB_WRITE_REPO_USER and/or \$GITLAB_WRITE_REPO_TOKEN CI variables must be set!"
  exit 1
fi

current_path=$PWD
project_path=$CI_PROJECT_DIR

# GitLab Pages branch name
branch=${1:-gitlab-pages}
# Path of the folder to push in GitLab Pages branch
source_path=${2:-$project_path/public}

repository_path=$project_path/$branch
# Path of the files to publish in GitLab Pages
pages_path=$repository_path/public
repository_url_path=$(echo $CI_PROJECT_URL | sed -E 's/^https?:\/\/(.+)$/\1/')
repository_url_protocol=$(echo $CI_PROJECT_URL | sed -E 's/^(https?:\/\/).+$/\1/')
repository_user=$GITLAB_WRITE_REPO_USER
repository_password=$GITLAB_WRITE_REPO_TOKEN
repository_url=${repository_url_protocol}${repository_user}:${repository_password}@${repository_url_path}
last_commit_message=$(git log -1 -z --format="%s%n%ncommit %H%nAuthor: %an <%ae>%nDate: %ad" HEAD)

if ! $(git clone --quiet --branch=$branch --single-branch $repository_url $repository_path); then
  # Create $branch if needed
  git clone --quiet --depth=1 --single-branch $repository_url $repository_path
  cd $repository_path
  git config user.name "$repository_user"
  git config user.email "${repository_user}@users.noreply.gitlab.com"
  git checkout --orphan $branch
  git rm --quiet -rf .
  cp -pa $project_path/.gitlab-ci.yml .gitlab-ci.yml

  echo "Git branch dedicated to [GitLab Pages](https://docs.gitlab.com/ce/user/project/pages/).

DO NOT EDIT THIS BRANCH BY HAND PLEASE" > README.md
  git add README.md .gitlab-ci.yml
  git commit -m "Git branch dedicated to GitLab Pages (https://docs.gitlab.com/ce/user/project/pages/)."
fi

cd $current_path

if [ -d $source_path ]; then
  cp -Rpa $source_path/. $pages_path
  cd $repository_path
  git config user.name "$repository_user"
  git config user.email "${repository_user}@users.noreply.gitlab.com"
  # Update git '$branch' and push changes to the current repository
  git add --all
  git diff-index --quiet HEAD || git commit -m "$last_commit_message"
  git push --quiet origin $branch | true
fi

cd $current_path
