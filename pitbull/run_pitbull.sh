#!/bin/bash

# This script downloads given passlist file (into passlist.txt) and runs btcrecover with it.
# It's the main functionality of Pitbull tool.
# It will run as a foreground process, and output progress to TTY. Some additional logs
# may be redirected to stderr (warnings, errors).
# By having the actual run call in separated file, we can easily modify the way
# Pitbull will be run (for example in detached TTY created with tmux).

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

source ./variables.sh

while getopts f:u:w: flag
do
    case "${flag}" in
        f) passlistFile=${OPTARG};;
        u) passlistFileUrl=${OPTARG};;
        w) walletString=${OPTARG};;
    esac
done

if [[ -z passlistFile ]] && [[ -z passlistFileUrl ]]; then
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
if [ -e "./$pipe" ]; then
  rm "./$pipe"
fi

mkfifo "./$pipe"

./capture_output.sh &

if [[ -z $passlistFile ]]; then
  script -f -c "./download_passlist.sh $passlistFileUrl $defaultPasslistFile && \
    ./run_btcrecover.sh $walletString $defaultPasslistFile" \
    $pipe
else
  script -f -c "./run_btcrecover.sh $walletString $passlistFile" $pipe
fi
