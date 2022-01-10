#!/bin/bash

# This is the main script of Pitbull tool. It's responsible for reading arguments
# and running the scripts in proper mode.
# By default, Pitbull will start new terminal session with tmux, under the name "pitbull".
# 
# Output will be redirected to "progress_view.txt" file via named pipe (btcrecover_out), but you can always
# re-attach terminal session with `tmux a -t "pitbull"` to see the live progress.

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

# Input args validation.
if [[ -z $fileUrl ]]; then
  echo "Passlist source missing"
  exit 1
fi

if [[ -z $walletString ]]; then
  echo "Wallet string missing"
  exit 1
fi

tmux new-session -d -s "pitbull" "./run_pitbull.sh $fileUrl $passlistFileName $walletString"
