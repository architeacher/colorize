#!/usr/bin/env bash

set -o pipefail

. "$(dirname "${0}")/common.sh"

exit_status=0

is_valid_combination() {
  local os="${1}" \
        arch="${2}"

  case "${os}" in
    "darwin" | "linux")
      case "${arch}" in
        "386" | "amd64" | "arm" | "arm64")
          return 0
          ;;
      esac
      ;;
    "windows")
      case "${arch}" in
        "386" | "amd64")
          return 0
          ;;
      esac
      ;;
  esac

  return 1
}

get_extension() {
  local target_os="${1}"

  [ "${target_os}" = "windows" ] && echo ".exe" || echo ""
}

build_binaries() {
  local src_dir="${1}" \
        output_dir="${2%/}" \
        target_os="${3}" \
        target_arch="${4}" \
        has_debug_flag="${5}" \
        build_branch \
        build_time \
        build_version \
        commit_sha \
        go_version

  build_branch="$(git symbolic-ref --short HEAD || git rev-parse --abbrev-ref HEAD &>/dev/null)"
  build_time="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

  build_version="dev"
  if [ -f "${src_dir}/../.version" ]; then
    build_version="$(cat "${src_dir}/../.version")"
  elif git describe --tags --exact-match &>/dev/null || git diff-index --quiet HEAD -- &>/dev/null; then
    build_version="$(git describe --always --abbrev=8 --tags --dirty='-Changes' 2> /dev/null | cut -d "v" -f 2 2>/dev/null)"
  fi

  commit_sha="$(git rev-list -1 HEAD &>/dev/null)"
  [ -n "$(git status --porcelain=v1 --untracked-files=no 2>/dev/null)" ] && commit_sha+="+CHANGES" 2>/dev/null

  go_version="$(go version | awk '{print $3}')"
  go_tags="netgo"

  local package_dir \
        package_name
  while IFS= read -r -d '' go_file; do
    if [ ! -f "${go_file}" ]; then
      log_warning "Failed to build ${go_file}"
      continue
    fi

    package_dir="$(dirname "${go_file}")"
    package_name="${package_dir##*/}"

    local binary_name
    binary_name="${output_dir}/${package_name}$(get_extension "${target_os}")"

    local linker_flags="-s -v -w"
    linker_flags+=" -X ${package_name}/build.Branch=${build_branch}"
    linker_flags+=" -X ${package_name}/build.Time=${build_time}"
    linker_flags+=" -X ${package_name}/build.Version=${build_version}"
    linker_flags+=" -X ${package_name}/build.CommitSHA=${commit_sha}"
    linker_flags+=" -X ${package_name}/build.GoVersion=${go_version}"

  local gc_flags=""
  if [[ "${has_debug_flag}" = "true" ]]; then
    log_info "Enabling debugging flags."
    gc_flags="all=-N -l"
  fi

    local build_command="CGO_ENABLED=0 GOOS=${target_os} GOARCH=${target_arch}"
    build_command+=" go build -v -o ${binary_name}"
    build_command+=" -installsuffix \"${target_os}.${target_arch}\""
    build_command+=" -asmflags='-trimpath=${package_dir}'"
    build_command+=" -gcflags \"${gc_flags}\""
    build_command+=" -ldflags \"${linker_flags}\""
    build_command+=" -tags ${go_tags}"
    build_command+=" -trimpath ${go_file}"

    log_info "Executing: ${build_command}"

    if ! eval "${build_command}"; then
      log_error "Failed to build ${go_file}"
      return 1
    fi

    log_info "Built ${binary_name}"
  done < <(find "${src_dir}" -type f -name "main.go" -print0)
}

display_help() {
  echo "Usage: ${0} [options]"
  echo
  echo "Options:"
  echo "  -s, --src-dir       Source directory containing Go files"
  echo "  -o, --output-dir    Output directory for the binaries"
  echo "  -os, --os           Target operating system (e.g., darwin, linux, windows)"
  echo "  -a, --arch          Target architecture (e.g., 386, amd64, arm, arm64)"
  echo "  -l, --log-level     Log level for messages (debug, info, notice, warning, error, critical, alert, emergency)"
  echo "  -d, --debug         Enable debugging flags"
  echo "  -h, --help          Display this help information"
  echo "  -x                  Display this execution flow"
  echo
}

parse_arguments() {
  local src_dir="" \
        output_dir="" \
        target_arch \
        target_os \
        log_level="info" \
        has_debug_flag=0

  # Set default target architecture and operating system
  target_arch="$(go env GOARCH)"
  target_os="$(go env GOOS)"

  while [[ "${#}" -gt 0 ]]; do
    case "${1}" in
      -s|--src-dir)
        src_dir="${2}"
        ;;
      -o|--output-dir)
        output_dir="${2}"
        ;;
      -os|--os)
        target_os="${2}"
        ;;
      -a|--arch)
        target_arch="${2}"
        ;;
      -l|--log-level)
        log_level="${2}"
        ;;
      -d|--debug)
        has_debug_flag="true"
        shift
        continue
        ;;
      -h|--help)
        display_help
        exit "${exit_status}"
        ;;
      -x)
        set -x
        shift
        continue
        ;;
      *)
        log_error "Unknown option: ${1}"
        display_help
        exit 1
        ;;
    esac
    shift
    shift
  done

  printf "%s\n%s\n%s\n%s\n%s\n%s" \
          "${src_dir}" \
          "${output_dir}" \
          "${target_os}" \
          "${target_arch}" \
          "${log_level}" \
          "${has_debug_flag}"
}

validate_arguments() {
  local src_dir="${1}" \
        output_dir="${2}" \
        target_os="${3}" \
        target_arch="${4}"

  if [ -z "${src_dir}" ] || [ -z "${output_dir}" ]; then
    abort "Missing required options: --src-dir and --output-dir are mandatory."
  fi

  if [[ -z "${target_os}" ]]; then
      log_error "Target operating system is not specified."
      display_help
  fi

  if [[ -z "${target_arch}" ]]; then
      log_error "Target architecture is not specified."
      display_help
  fi

  if ! is_valid_combination "${target_os}" "${target_arch}"; then
    abort "Invalid combination: ${target_os}/${target_arch}"
  fi
}

script_cleanup() {
 local output_dir="${1}"

 log_info "cleaning up"

 rm -rf "${output_dir}"
}

main() {
  trap 'exit ${exit_status}' EXIT

  validate_bash

  {
    read -r src_dir
    read -r output_dir
    read -r target_os
    read -r target_arch
    read -r log_level
    read -r has_debug_flag
  } <<< "$(parse_arguments "${@}")"

  validate_arguments "${src_dir}" "${output_dir}" "${target_os}" "${target_arch}"

  log_set_level "$(parse_log_level "${log_level}")"

  build_binaries "${src_dir}" "${output_dir}" "${target_os}" "${target_arch}" "${has_debug_flag}"
}

main "${@}"
