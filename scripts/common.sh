#!/usr/bin/env bash

set -eou pipefail

. "$(dirname "${0}")/logger.sh"

SUDO=

is_command() {
    command -v "${1}" 2> /dev/null
}

check_sudo() {
  [ "$(id -u)" -eq 0 ] && return

  SUDO=$(is_command "sudo")
  [ ! -x "${SUDO}" ] && abort "This script must be executed as root."
}

get_profile() {
  local shell_profile="${HOME}/.bash_profile"

  case "${SHELL}" in
    */bash*)
      [ -r "${HOME}/.bash_profile" ] && shell_profile="${HOME}/.bash_profile" || shell_profile="${HOME}/.profile"
      ;;
    */zsh*)
      shell_profile="${HOME}/.zprofile"
      ;;
    *)
      shell_profile="${HOME}/.profile"
      ;;
  esac

  print "${shell_profile}"
}

get_random_string() {
  print "$(export LC_CTYPE=C; cat </dev/urandom | tr -dc 'a-zA-Z0-9\.' | fold -w 32 | head -n 1)"
}

validate_bash() {
  # Fail fast with a concise message when not using bash
  # Single brackets are needed here for POSIX compatibility
  # shellcheck disable=SC2292
  [ -z "${BASH_VERSION:-}" ] && abort "Bash is required to interpret this script."

  # Check if running in a compatible bash version.
  ((BASH_VERSINFO[0] < 3)) && abort "Bash version 3 or above is required."

  # Check if script is run in POSIX mode.
  if [[ -n "${POSIXLY_CORRECT+1}" ]]
  then
    abort "Bash must not run in POSIX compatibility mode. Please disable by unsetting POSIXLY_CORRECT and try again."
  fi
}
