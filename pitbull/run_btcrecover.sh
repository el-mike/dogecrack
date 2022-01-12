#!/bin/bash

# This script is responsible for running btcrecover process.

walletString=$1
passlistFileName=$2

python3 ./btcrecover/btcrecover.py --dsw \
  --data-extract-string $walletString \
  --passwordlist $passlistFileName \
  --enable-gpu
