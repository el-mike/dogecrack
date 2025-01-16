#!/bin/bash

docker container prune -f

docker build -t michalhuras/pitbull:latest .

docker login
docker push michalhuras/pitbull:latest
