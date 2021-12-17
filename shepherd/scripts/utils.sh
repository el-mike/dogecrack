#!/bin/bash

get_fake_vast_container_id() {
  echo $(docker ps | grep "fake_vast" | awk '{ print $1 }')
}

get_fake_vast_container_ip() {
  id=$(get_fake_vast_container_id)
  echo $(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $id)
}


