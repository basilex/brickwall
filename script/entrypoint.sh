#!/usr/bin/env bash
#
# Abort on any error (including if wait-for-it fails).
#
set -e
#
# Wait for the backend to be up, if we know where it is.
#
if [ -n "$DBS_HOST" ]; then
    /app/wait-for-it.sh "${DBS_HOST}:${DBS_PORT:-5432}"
fi
if [ -n "$BUS_HOST" ]; then
    /app/wait-for-it.sh "${BUS_HOST}:${BUS_PORT:-4222}"
fi
#
# Run the main container command.
#
exec "$@"
