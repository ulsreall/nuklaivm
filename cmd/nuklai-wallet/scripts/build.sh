#!/usr/bin/env bash
# Copyright (C) 2024, AllianceBlock. All rights reserved.
# See the file LICENSE for licensing terms.

set -o errexit
set -o nounset
set -o pipefail

# Set the CGO flags to use the portable version of BLST
#
# We use "export" here instead of just setting a bash variable because we need
# to pass this flag to all child processes spawned by the shell.
export CGO_CFLAGS="-O -D__BLST_PORTABLE__"

PUBLISH=${PUBLISH:-true}

# Install wails
go install -v github.com/wailsapp/wails/v2/cmd/wails@v2.8.0

# Build file for local arch
#
# Don't use upx: https://github.com/upx/upx/issues/446
wails build -clean -platform darwin/universal

OUTPUT=build/bin/Nuklai\ Wallet.app
if [ ! -d "$OUTPUT" ]; then
  exit 1
fi

# Remove any previous build artifacts
rm -rf nuklai-wallet.zip

# Exit early if not publishing
if [ "${PUBLISH}" == false ]; then
  echo "not publishing app"
  ditto -c -k --keepParent build/bin/Nuklai\ Wallet.app nuklai-wallet.zip
  exit 0
fi
echo "publishing app"

# Sign code
codesign -s "${APP_SIGNING_KEY_ID}" --deep  --timestamp -o runtime -v build/bin/Nuklai\ Wallet.app
ditto -c -k --keepParent build/bin/Nuklai\ Wallet.app nuklai-wallet.zip

# Need to sign to allow for app to be opened on other computers
xcrun altool --notarize-app --primary-bundle-id ai.nukl.nuklai-wallet --username "${APPLE_NOTARIZATION_USERNAME}" --password "@keychain:altool" --file nuklai-wallet.zip

# Log until exit
read -p -r "Input APPLE_NOTARIZATION_REQUEST_ID: " APPLE_NOTARIZATION_REQUEST_ID
while true
do
  xcrun altool --notarization-info "${APPLE_NOTARIZATION_REQUEST_ID}" -u "${APPLE_NOTARIZATION_USERNAME}" -p "@keychain:altool"
  sleep 15
done
