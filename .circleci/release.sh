#!/usr/bin/env bash

set -o errexit
set -o nounset
set -euo pipefail

# curl -sL https://git.io/goreleaser | bash
mkdir -p "$HOME/.ssh/" && echo -e "Host github.com\n\tStrictHostKeyChecking no\n" > ~/.ssh/config

VERSION=$(cat ./.version)
LAST_COMMIT_MESSAGE=$(git log -1 --pretty=%s)

if [ -n "$VERSION" ]; then
  echo "Create release: ${VERSION}"
  curl -S "${GITHUB_API}/repos/${REPO_NAME}/${IMAGE_NAME}/releases" \
       -H "Accept: application/vnd.github.v3+json"                  \
       -H "Authorization: token ${GITHUB_ACCESS_TOKEN}"             \
       -H "Content-Type: application/json; charset=UTF-8"           \
       -d @- << REQUEST_BODY
          {
            "name": "${VERSION}",
            "tag_name": "${VERSION}",
            "body": "${LAST_COMMIT_MESSAGE}"
          }
REQUEST_BODY
fi
