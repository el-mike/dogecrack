#!/bin/bash

source $(dirname "$0")/utils.sh

echo $(get_fake_vast_container_ip $1)
