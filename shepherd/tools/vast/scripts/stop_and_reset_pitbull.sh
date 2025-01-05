#!/bin/bash

source $(dirname "$0")/utils.sh

id_one=$(get_fake_vast_container_id 1)
id_two=$(get_fake_vast_container_id 2)

# For some reason, "docker exec -it" does not see /app in $PATH
docker exec -it $id_one bash -c "/app/pitbull.sh kill ; rm -f /app/progress_view.txt"
docker exec -it $id_two bash -c "/app/pitbull.sh kill ; rm -f /app/progress_view.txt"
