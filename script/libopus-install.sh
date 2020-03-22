#!/usr/bin/env sh

OS=$(uname)
case "$OS" in
  'Linux') sudo apt-get install pkg-config libopus-dev libopusfile-dev ;; # Linux
  'Darwin') brew install pkg-config opus opusfile ;; # Mac OSX
  'WindowsNT') choco install pkgconfiglite opus-tools ;; # Windows
  *)        printf "unknown: %s\n" "$OS"; exit 1 ;;
esac
