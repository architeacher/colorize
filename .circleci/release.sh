#!/usr/bin/env bash

set -o errexit
set -o nounset
set -euo pipefail

# curl -sL https://git.io/goreleaser | bash
mkdir -p "$HOME/.ssh/" && echo -e "Host github.com\n\tStrictHostKeyChecking no\n" > ~/.ssh/config

VERSION=$(cat ./.version)

if [ -n "$VERSION" ]; then
  echo "Creating release: ${VERSION}"
  curl -X POST                                                  \
       -H "Accept: application/vnd.github.v3+json"              \
       -H "Authorization: token ${GITHUB_ACCESS_TOKEN}"         \
       -d "{\"tag_name\": \"${VERSION}\"}"                      \
       "${GITHUB_API}repos/${REPO_NAME}/${IMAGE_NAME}/releases"
fi
