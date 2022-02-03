#!/bin/bash

# $1 - host IP
# $2 - ssh config file path

figerprintExists=$(ssh-keygen -F $1)

if [[ -z $figerprintExists ]]
then
  echo "Adding fingerprint for host $1"
  ssh-keyscan -H $1 >> $2
else
  echo "Host already added"
fi

