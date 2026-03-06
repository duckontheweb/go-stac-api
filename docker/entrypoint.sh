#! /bin/bash

if [[ -v DEBUG ]]; then
    debug=$(echo "$DEBUG" | tr '[:upper:]' '[:lower:]')
    if [[ "$debug" = "yes" ]] || [[ "$debug" = "true" ]]; then
        config_file=".air.debug.toml"
    else
        config_file=".air.toml"
    fi
else
    config_file=".air.toml"
fi

air -c "${config_file}"
