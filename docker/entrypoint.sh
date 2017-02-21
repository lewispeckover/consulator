#!/bin/dumb-init /bin/sh
set -e

if [ "$(basename $1 2>/dev/null)" != 'consulator' ]; then
    set -- consulator "$@"
fi

exec "$@"
