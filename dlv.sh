#!/bin/sh
cd /go/src/work
if [ "$(echo "$DEBUG_MODE" | tr '[:upper:]' '[:lower:]')" = "true" ]; then
  dlv debug --headless --log -l 0.0.0.0:2345 --api-version=2
else
  air
fi