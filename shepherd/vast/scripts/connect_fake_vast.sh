#!/bin/bash

source $(dirname "$0")/utils.sh

instanceNumber=$1

ip=$(get_fake_vast_container_ip $1)

# ssh complains about re-connecting to this host after rebuilding an image,
# therefore we remove this entry.
# TODO: investigate.
ssh-keygen -f "$HOME/.ssh/known_hosts" -R "$ip"

ssh root@$ip
