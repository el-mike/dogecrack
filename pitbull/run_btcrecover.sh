#!/bin/bash

# This script is responsible for running btcrecover process.

dirname=$(dirname "$0")

walletString=$1
passlistFileName=$2

python3 $dirname/btcrecover/btcrecover.py --dsw \
  --data-extract-string $walletString \
  --passwordlist $passlistFileName \
  --enable-gpu
