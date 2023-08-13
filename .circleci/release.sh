#!/usr/bin/env bash

set -o errexit
set -o nounset
set -euo pipefail

# curl -sL https://git.io/goreleaser | bash

version="$(cat ./.versoin)"
last_commit_message="$(git log -1 --pretty=%s)"

if [ -n "${version}" ]; then
  echo "Create release: ${version}"
  curl -S "${GITHUB_API}/repos/${REPO_NAME}/${IMAGE_NAME}/releases" \
       -H "Accept: application/vnd.github.v3+json" \
       -H "Authorization: token ${GITHUB_ACCESS_TOKEN}" \
       -H "Content-Type: application/json; charset=UTF-8" \
       -d @- << REQUEST_BODY
          {
            "name": "${version}",
            "tag_name": "${version}",
            "body": "${last_commit_message}"
          }
REQUEST_BODY
fi
