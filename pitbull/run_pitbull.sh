#!/bin/bash

# This script downloads given passlist file (into passlist.txt) and runs btcrevcover with it.
# It's the main functionality of Pitbull tool.
# It will run as a foreground process, and output progress to TTY. Some additional logs
# may be redirected to stderr (warnings, errors).
# By having the actual run call in separated file, we can easily modify the way
# Pitbull will be run (for example in detached TTY created with tmux).

fileUrl=$1
passlistFileName=$2
walletString=$3

pipe='btcrecover_out'

# Output capture setup.
# If pipe exists, remove it - it ensures that no other agent is
# reading from the output pipe.
# For some reason, "-p" (testing for named pipe exactly) does not work sometimes,
# therefore we use "-e" instead. 
if [[ -e $pipe ]]; then
  rm "$pipe"
fi

mkfifo "$pipe" && ./capture_output.sh &

script -f -c "./download_passlist.sh $fileUrl $passlistFileName && \
  ./run_btcrecover.sh $walletString $passlistFileName" \
  $pipe
