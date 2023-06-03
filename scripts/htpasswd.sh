#!/bin/sh
set -e

readonly user="$1"

echo $(htpasswd -nB $user) | sed -e s/\\$/\\$\\$/g