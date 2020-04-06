#!/usr/bin/env bash
#

# === gorelease, see: https://goreleaser.com/ ===
goreleaser --snapshot --skip-publish --rm-dist
cd bin
python3 -m http.server