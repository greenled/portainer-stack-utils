#!/usr/bin/env bash
set -e
[[ "$TRACE" ]] && set -x

bash "$(dirname "$0")/../tests/run.sh"
