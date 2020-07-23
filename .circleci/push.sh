#!/usr/bin/env bash

set -o errexit
set -o nounset
set -euo pipefail

mkdir -p "$HOME/.ssh/" && echo -e "Host github.com\n\tStrictHostKeyChecking no\n" > ~/.ssh/config

VERSION=$(cat .version 2> /dev/null)

if [ -n "$VERSION" ]; then
  git tag "$VERSION"
  echo "$VERSION -> $(git rev-parse --short=8 "$VERSION" 2> /dev/null)"
  git push "https://${GITHUB_ACCESS_TOKEN}:x-oauth-basic@github.com/${REPO_NAME}/${IMAGE_NAME}" --tags
fi
