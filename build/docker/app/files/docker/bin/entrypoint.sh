#!/usr/bin/env sh

if [ -n "$SLEEP_BEFORE_START" ]; then
    sleep $SLEEP_BEFORE_START
fi

exec ./bin "$@"