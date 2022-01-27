#!/bin/bash

# Runs actual btcrecover process.

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

walletString=$1
passlistFile=$2

python3 ./btcrecover/btcrecover.py --dsw \
  --data-extract-string $walletString \
  --passwordlist $passlistFile \
  --enable-gpu
