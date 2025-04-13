#!/bin/sh
set -e

# Make sure the config directory exists
mkdir -p configs

# Write the config file
echo "$CONFIG_YML" > configs/config.yaml
echo "$CONFIG_YML" > configs/config.local.yaml

# Run the Go binary
./order-management