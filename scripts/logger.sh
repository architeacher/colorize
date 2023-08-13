#!/usr/bin/env bash

set -eou pipefail

readonly LOG_LEVEL_DEBUG=7
readonly LOG_LEVEL_INFO=6
readonly LOG_LEVEL_NOTICE=5
readonly LOG_LEVEL_WARNING=4
readonly LOG_LEVEL_ERROR=3
readonly LOG_LEVEL_CRITICAL=2
readonly LOG_LEVEL_ALERT=1
readonly LOG_LEVEL_EMERGENCY=0

_log_level="${LOG_LEVEL_INFO}"

# string formatters, True if go_file descriptor FD is open and refers to a terminal.
tty_escape() { :; }

if [[ -t 1 ]]
then
  tty_escape() { printf "\x1b[%sm" "${1}"; }
fi

tty_4bit_mk() { tty_escape "${1}"; }
tty_8bit_mk() { tty_escape "38;5;${1}"; }
tty_4bit_mkbold() { tty_4bit_mk "${1};1"; }
tty_8bit_mkbold() { tty_8bit_mk "${1};1"; }

tty_bold="$(tty_4bit_mkbold 39)"
tty_underline="$(tty_escape "4;39")"
tty_reset="$(tty_escape 0)"

tty_blue="$(tty_4bit_mk 34)"
tty_corn="$(tty_8bit_mk 178)"
tty_cyan="$(tty_8bit_mk 45)"
tty_green="$(tty_8bit_mk 35)"
tty_imperial="$(tty_4bit_mk 91)"
tty_lime="$(tty_8bit_mk 113)"
tty_magenta="$(tty_8bit_mk 170)"
tty_move="$(tty_4bit_mk 35)"
tty_olive="$(tty_8bit_mk 64)"
tty_orange="$(tty_8bit_mk 208)"
tty_pink="$(tty_8bit_mk 198)"
tty_red="$(tty_8bit_mkbold 196)"
tty_teal="$(tty_4bit_mk 36)"
tty_yellow="$(tty_8bit_mk 227)"

log_set_level() {
	_log_level="$1"
}

log_priority() {
  local log_level="${1}"

  if test -z "${log_level}"; then
    echo "${log_level}"
    return
  fi

  [ "${log_level}" -le "${_log_level}" ]
}

display_message() {
  echo -e "${@}" 1>&2
}

log_color() {
  local log_level="${1}" \
        tty_color=""

  case "${log_level}" in
    "${LOG_LEVEL_DEBUG}")
      tty_color="${tty_cyan}"
      ;;
    "${LOG_LEVEL_INFO}")
      tty_color="${tty_olive}"
      ;;
    "${LOG_LEVEL_NOTICE}")
      tty_color="${tty_green}"
      ;;
    "${LOG_LEVEL_WARNING}")
      tty_color="${tty_corn}"
      ;;
    "${LOG_LEVEL_ERROR}")
      tty_color="${tty_magenta}"
      ;;
    "${LOG_LEVEL_CRITICAL}")
      tty_color="${tty_imperial}"
      ;;
    "${LOG_LEVEL_ALERT}")
      tty_color="${tty_pink}"
      ;;
    "${LOG_LEVEL_EMERGENCY}")
      tty_color="${tty_red}"
      ;;
    *)
      tty_color="${tty_bold}"
      ;;
   esac

  printf "%s" "${tty_color}"
}

log_tag() {
  local log_level="${1}" \
        tag=""

  case "${log_level}" in
    "${LOG_LEVEL_DEBUG}")
      tag="Debug"
      ;;
    "${LOG_LEVEL_INFO}")
      tag="Info"
      ;;
    "${LOG_LEVEL_NOTICE}")
      tag="Notice"
      ;;
    "${LOG_LEVEL_WARNING}")
      tag="Warning"
      ;;
    "${LOG_LEVEL_ERROR}")
      tag="Error"
      ;;
    "${LOG_LEVEL_CRITICAL}")
      tag="Critical"
      ;;
    "${LOG_LEVEL_ALERT}")
      tag="Alert"
      ;;
    "${LOG_LEVEL_EMERGENCY}")
      tag="Emergency"
      ;;
    *)
      tag=""
      ;;
  esac

  printf "%s" "${tag}"
}

log_prefix() {
  echo "==>"
}

log() {
  local log_level="${1}" \
        message="" \
        tty_color="" \
        tag=""

  {
    read -r tty_color
  } <<< "$(log_color "${log_level}")"

  {
    read -r tag
  } <<< "$(log_tag "${log_level}")"

  shift
  message="${*}"

  if [ "${tty_color}" != "" ]; then
    display_message "${tty_color}" "$(log_prefix)" "${tag}: " "${message}" "${tty_reset}"
  else
    printf "%s %s" "$(log_prefix)" "${message}"
  fi
}

log_debug() {
  log_priority "${LOG_LEVEL_DEBUG}" || return 0
  log "${LOG_LEVEL_DEBUG}" "${@}"
}

log_info() {
  log_priority "${LOG_LEVEL_INFO}" || return 0
  log "${LOG_LEVEL_INFO}" "${@}"
}

log_notice() {
  log_priority "${LOG_LEVEL_NOTICE}" || return 0
  log "${LOG_LEVEL_NOTICE}" "${@}"
}

log_warning() {
  log_priority "${LOG_LEVEL_WARNING}" || return 0
  log "${LOG_LEVEL_WARNING}" "${@}"
}

log_error() {
  log_priority "${LOG_LEVEL_ERROR}" || return 0
  log "${LOG_LEVEL_ERROR}" "${@}"
}

log_critical() {
  log_priority "${LOG_LEVEL_CRITICAL}" || return 0
  log "${LOG_LEVEL_CRITICAL}" "${@}"
}

log_alert() {
  log_priority "${LOG_LEVEL_ALERT}" || return 0
  log "${LOG_LEVEL_ALERT}" "${@}"
}

log_emergency() {
  log_priority "${LOG_LEVEL_EMERGENCY}" || return 0
  log "${LOG_LEVEL_EMERGENCY}" "${@}"
}

# abort with an error message.
abort() {
  local message="${1}"

  log_error "${message}"
  exit 1
}

parse_log_level() {
  local log_input="${1}" \
        log_level=""

  case "${log_input}" in
      "debug")
        log_level="${LOG_LEVEL_DEBUG}"
        ;;
      "info")
        log_level="${LOG_LEVEL_INFO}"
        ;;
      "notice")
        log_level="${LOG_LEVEL_NOTICE}"
        ;;
      "warning")
        log_level="${LOG_LEVEL_WARNING}"
        ;;
      "error")
        log_level="${LOG_LEVEL_ERROR}"
        ;;
      "critical")
        log_level="${LOG_LEVEL_CRITICAL}"
        ;;
      "alert")
        log_level="${LOG_LEVEL_ALERT}"
        ;;
      "emergency")
        log_level="${LOG_LEVEL_EMERGENCY}"
        ;;
      *)
        abort "Invalid log level: ${log_level}"
        ;;
    esac

    echo "${log_level}"
}
