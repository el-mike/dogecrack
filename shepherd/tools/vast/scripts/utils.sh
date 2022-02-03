#!/bin/bash

# $1 contains the number of local fake vast instance.
get_fake_vast_container_id() {
  echo $(docker ps | grep "[f]ake_vast_$1" | awk '{ print $1 }')
}

# $1 contains the number of local fake vast instance.
get_fake_vast_container_ip() {
  id=$(get_fake_vast_container_id $1)
  echo $(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $id)
}


