#!/bin/bash

# This script downloads given passlist file (into passlist.txt) and runs btcrevcover with it.
# It's the main functionality of Pitbull tool.
# It will run as a foreground process, and output progress to TTY. Some additional logs
# may be redirected to stderr (warnings, errors).

dirname=$(dirname "$0")

passlistFileName='passlist.txt'
pipe='btcrecover_out'

while getopts f:w: flag
do
    case "${flag}" in
        f) fileUrl=${OPTARG};;
        w) walletString=${OPTARG};;
    esac
done

echo "Wallet string: $walletString"

# Input args validation.
if [[ -z $fileUrl ]]; then
  echo "Passlist source missing"
  exit 1
fi

if [[ -z $walletString ]]; then
  echo "Wallet string missing"
  exit 1
fi

# Output capture setup.
# If pipe exists, remove it - it ensures that no other agent is
# reading from the output pipe.
# For some reason, "-p" (testing for named pipe exactly) does not work sometimes,
# therefore we use "-e" instead. 
if [[ -e $pipe ]]; then
  rm "$pipe"
fi

mkfifo "$pipe" && ./capture_output.sh &

script -f -c "$dirname/download_passlist.sh $fileUrl $passlistFileName && \
  $dirname/run_btcrecover.sh $walletString $passlistFileName" \
  $pipe
