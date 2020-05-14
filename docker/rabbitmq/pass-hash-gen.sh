#!/usr/bin/env bash

# Use this to generate a password hash for the definitions.json 
# file. 
# 
# More info here: https://www.rabbitmq.com/passwords.html
# 
# Only works when 
#   "hashing_algorithm": "rabbit_password_hashing_sha256",

set -o errexit
set -o nounset

declare -r passwd="${1:-newpassword}"

declare -r tmp0="$(mktemp)"
declare -r tmp1="$(mktemp)"

function onexit
{
    rm -f "$tmp0"
    rm -f "$tmp1"
}

trap onexit EXIT

dd if=/dev/urandom of="$tmp0" count=4 bs=1 > /dev/null 2>&1
cp -f "$tmp0" "$tmp1"
echo -n "$passwd" >> "$tmp0"
openssl dgst -binary -sha256 "$tmp0" >> "$tmp1"
openssl base64 -A -in "$tmp1"