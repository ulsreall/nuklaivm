#!/usr/bin/env bash
# Copyright (C) 2024, AllianceBlock. All rights reserved.
# See the file LICENSE for licensing terms.

set -e

# Set the CGO flags to use the portable version of BLST
#
# We use "export" here instead of just setting a bash variable because we need
# to pass this flag to all child processes spawned by the shell.
export CGO_CFLAGS="-O -D__BLST_PORTABLE__" CGO_ENABLED=1

if ! [[ "$0" =~ scripts/tests.unit.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# Use mapfile to read the output into an array
mapfile -t packages < <(go list ./... | grep -v tests)
# Now use the array with the 'go test' command
go test -race -timeout="10m" -coverprofile="coverage.out" -covermode="atomic" "${packages[@]}"
