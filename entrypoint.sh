#!/bin/sh
set -e

# Make sure the config directory exists
mkdir -p config

# Write the config file
echo "$CONFIG_YML" > config/config.yml

# Run the Go binary
./order-management