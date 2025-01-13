#!/bin/bash

docker build -t michalhuras/shepherd:latest -f cmd/shepherd/prod.Dockerfile .
