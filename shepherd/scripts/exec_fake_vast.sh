#!/bin/bash

source $(dirname "$0")/utils.sh

id=$(get_fake_vast_container_id)

docker exec -it $id bash
