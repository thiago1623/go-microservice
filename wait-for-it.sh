#!/bin/sh
# wait-for-it.sh: Script to wait for a service to be ready
# Usage: wait-for-it.sh host:port [-t timeout] [-- command args]
# https://github.com/vishnubob/wait-for-it

set -e

host_port="$1"
shift
timeout=30

while ! nc -z ${host_port} 2>/dev/null; do
  echo "Waiting for ${host_port}..."
  sleep 1
  timeout=$(expr ${timeout} - 1)
  if [ ${timeout} -eq 0 ]; then
    echo "Timeout occurred waiting for ${host_port}"
    exit 1
  fi
done

echo "${host_port} is ready!"
exec "$@"
