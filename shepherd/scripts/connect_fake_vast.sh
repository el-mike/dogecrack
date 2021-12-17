#!/bin/bash

source $(dirname "$0")/utils.sh

ip=$(get_fake_vast_container_ip)

# ssh complains about re-connecting to this host after rebuilding an image,
# therefore we remove this entry.
# TODO: investigate.
ssh-keygen -f "/home/elmike/.ssh/known_hosts" -R "$ip"

echo root@$ip

ssh root@$ip
