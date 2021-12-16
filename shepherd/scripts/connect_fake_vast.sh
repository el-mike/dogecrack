#!/bin/bash

id=$(docker ps | grep "fake_vast" | awk '{ print $1 }')
ip=$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $id)

# ssh complains about re-connecting to this host after rebuilding an image,
# therefore we remove this entry.
# TODO: investigate.
ssh-keygen -f "/home/elmike/.ssh/known_hosts" -R "$ip"

ssh root@$ip
