#!/usr/bin/env bash

# This script will install the ASDF plugins required for this project

set -euo pipefail
IFS=$'\n\t'

# Temporary transition over to mise from asdf
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

setup_mise() {
  temp_MISE_SHORTHANDS_FILE=$(mktemp)
  trap 'do_mise_install' EXIT

  do_mise_install() {
    cat "$temp_MISE_SHORTHANDS_FILE"
    MISE_SHORTHANDS_FILE=$temp_MISE_SHORTHANDS_FILE $MISE_COMMAND install
    rm -f "$temp_MISE_SHORTHANDS_FILE"
  }

  install_plugin() {
    local plugin=$1
    local source=${2-}

    # No source? mise defaults should suffice.
    if [[ -z $source ]]; then return; fi

    # See https://mise.jdx.dev/configuration.html#mise-shorthands-file-config-mise-shorthands-toml
    echo "$plugin = \"$source\"" >>"$temp_MISE_SHORTHANDS_FILE"
  }

  remove_plugin_with_source() {
    local plugin=$1
    local source=$2

    if ! $MISE_COMMAND plugin list --urls | grep -qF "${source}"; then
      return
    fi

    echo "# Removing plugin ${plugin} installed from ${source}"
    $MISE_COMMAND plugin remove "${plugin}" || {
      echo "Failed to remove plugin: ${plugin}"
      exit 1
    } >&2
  }
}

if command -v mise >/dev/null; then
  MISE_COMMAND=$(which mise)
  export MISE_COMMAND
  setup_mise
elif command -v rtx >/dev/null; then
  MISE_COMMAND=$(which rtx)
  export MISE_COMMAND
  setup_mise
elif [[ -n ${ASDF_DIR-} ]]; then
  setup_asdf
fi

install_plugin golang # Install golang first as some of the other plugins require it.
install_plugin goreleaser
install_plugin golangci-lint
install_plugin shfmt
install_plugin shellcheck
install_plugin pre-commit
