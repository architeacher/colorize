#!/usr/bin/env bash

set -o errexit
set -o nounset
set -euo pipefail

last_commit_message=$(git log -1 --pretty=%B)
version_file="./.version"
touch "${version_file}"

versionize() {
    local commit_message="${1}"
    local major="${2}"
    local minor="${3}"
    local patch="${4}"
    local version_file="${5}"

    if echo "${commit_message}" | grep -iqE "\[major\]"; then
      major=$((major+1))
      echo "v${major}.0.0" > "${version_file}"
    elif echo "${commit_message}" | grep -iqE "\[minor\]"; then
      minor=$((minor+1))
      echo "v${major}.${minor}.0" > "${version_file}"
    elif echo "${commit_message}" | grep -iqE "\[patch\]"; then
      patch=$((patch+1))
      echo "v${major}.${minor}.${patch}" > "${version_file}"
    fi
}

if version=$(git describe --abbrev=0 --tags 2> /dev/null) && [[ (-n "$(git diff "${version}")") || (-z "${version}") ]]; then
  version=${version:-'0.0.0'}; version=${version#"v"}
  major=${version%%.*}; version=${version#*.}
  minor=${version%%.*}; version=${version#*.}
  patch=${version%%.*}; version=${version#*.}

  versionize "${last_commit_message}" "${major}" "${minor}" "${patch}" "${version_file}"
else
  versionize "${last_commit_message}" 0 0 0 "${version_file}"
fi
