#!/bin/bash
set -e

# Bring schema up to date non-destructively before either service starts.
/app/wiki sync

/app/wiki serve &
BACK=$!
node /app/.output/server/index.mjs &
FRONT=$!

_stop() { kill "$BACK" "$FRONT" 2>/dev/null; exit; }
trap _stop INT TERM

# Exit (and let Docker restart the container) if either service dies.
wait -n "$BACK" "$FRONT"
_stop
