#!/bin/bash

source $(dirname "$0")/utils.sh

id=$(get_fake_vast_container_id $1)

docker exec -it $id bash
