#!/bin/bash

docker container prune -f

docker image rm michalhuras/pitbull:dev_local

docker build -t michalhuras/pitbull:dev_local .
