#!/usr/bin/env bash

# This script will install the ASDF plugins required for this project

set -euo pipefail
IFS=$'\n\t'

# shellcheck source=/dev/null
source "$ASDF_DIR/asdf.sh"

plugin_list=$(asdf plugin list || echo)

install_plugin() {
  local plugin=$1

  if ! echo "${plugin_list}" | grep -q "${plugin}"; then
    echo "# Installing plugin" "$@"
    asdf plugin add "$@" || {
      echo "Failed to install plugin:" "$@"
      exit 1
    } >&2
  fi

  echo "# Installing ${plugin} version"
  asdf install "${plugin}" || {
    echo "Failed to install plugin version: ${plugin}"
    exit 1
  } >&2

  # Use this plugin for the rest of the install-asdf-plugins.sh script...
  asdf shell "${plugin}" "$(asdf current "${plugin}" | awk '{print $2}')"
}

remove_plugin_with_source() {
  local plugin=$1
  local source=$2

  if ! asdf plugin list --urls | grep -qF "${source}"; then
    return
  fi

  echo "# Removing plugin ${plugin} installed from ${source}"
  asdf plugin remove "${plugin}" || {
    echo "Failed to remove plugin: ${plugin}"
    exit 1
  } >&2

  # Refresh list of installed plugins.
  plugin_list=$(asdf plugin list)
}

install_plugin golang
install_plugin goreleaser
install_plugin golangci-lint
install_plugin shfmt
install_plugin shellcheck
install_plugin pre-commit
