#!/bin/bash

# Base theme configuration
# Exports common variables if needed

init_theme() {
    export TRONCLI_THEME="$1"
    echo "TronCLI Theme set to: $TRONCLI_THEME"
}
