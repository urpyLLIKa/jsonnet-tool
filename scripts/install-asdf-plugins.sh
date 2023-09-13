#!/usr/bin/env bash

# This script will install the ASDF plugins required for this project

set -euo pipefail
IFS=$'\n\t'

# Temporary transition over to rtx from asdf
# see https://gitlab.com/gitlab-com/runbooks/-/issues/134
# for details
setup_asdf() {
  # shellcheck source=/dev/null
  source "$ASDF_DIR/asdf.sh"

  plugin_list=$(asdf plugin list || echo "")

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
}

setup_rtx() {
  temp_RTX_SHORTHANDS_FILE=$(mktemp)
  trap 'do_rtx_install' EXIT

  do_rtx_install() {
    cat "$temp_RTX_SHORTHANDS_FILE"
    RTX_SHORTHANDS_FILE=$temp_RTX_SHORTHANDS_FILE rtx install
    rm -f "$temp_RTX_SHORTHANDS_FILE"
  }

  install_plugin() {
    local plugin=$1
    local source=${2-}

    # No source? rtx defaults should suffice.
    if [[ -z $source ]]; then return; fi

    # See https://github.com/jdxcode/rtx#rtx_shorthands_fileconfigrtxshorthandstoml
    echo "$plugin = \"$source\"" >>"$temp_RTX_SHORTHANDS_FILE"
  }

  remove_plugin_with_source() {
    local plugin=$1
    local source=$2

    if ! rtx plugin list --urls | grep -qF "${source}"; then
      return
    fi

    echo "# Removing plugin ${plugin} installed from ${source}"
    rtx plugin remove "${plugin}" || {
      echo "Failed to remove plugin: ${plugin}"
      exit 1
    } >&2
  }
}

if command -v rtx >/dev/null; then
  setup_rtx
elif [[ -n ${ASDF_DIR-} ]]; then
  setup_asdf
fi

install_plugin golang # Install golang first as some of the other plugins require it.
install_plugin goreleaser
install_plugin golangci-lint
install_plugin shfmt
install_plugin shellcheck
install_plugin pre-commit
