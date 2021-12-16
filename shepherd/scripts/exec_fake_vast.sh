#!/bin/bash

id=$(docker ps | grep "fake_vast" | awk '{ print $1 }')

docker exec -it $id bash
