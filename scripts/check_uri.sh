#!/bin/sh
set -e

readonly method="$1"
readonly uri="$2"

check_server() {
    http_code=$(curl --silent --output /dev/null --write-out "%{http_code}" -X "$method" "$uri")
    if [ $http_code -ge 200 ] && [ $http_code -lt 500 ]; then
        return 0
    fi
    return 1
}

while ! check_server; do
    echo "Ping $method $uri failed."
    echo "Waiting for server. Retrying in 10 seconds..."
    sleep 10
done