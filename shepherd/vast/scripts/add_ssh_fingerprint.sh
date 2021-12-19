#!/bin/bash

# $1 - host IP
# $2 - ssh config file path

echo $1
echo $2

ssh-keyscan -H $1 >> $2
