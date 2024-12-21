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

while getopts t:f:u:w: flag
do
    case "${flag}" in
        t) tokenList=${OPTARG};;
        f) passlistFile=${OPTARG};;
        u) passlistFileUrl=${OPTARG};;
        w) walletString=${OPTARG};;
    esac
done

if [[ -z $walletString ]]; then
  echo "Wallet string missing"
  exit 1
fi

if [[ -z tokenList ]] && [[ -z passlistFile ]] && [[ -z passlistFileUrl ]]; then
  echo "Passlist source missing"
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

echo $tokenList

if [[ -n $tokenList ]]; then
  # We need to use -e flag to preserve newline characters from tokenlist argument.
  echo -e "$tokenList" > "$defaultTokenlistFile"

  script -f -c "python3 ./btcrecover/btcrecover.py --dsw --data-extract-string $walletString \
    --tokenlist $defaultTokenlistFile --enable-gpu" \
    $pipe
elif [[ -n $passlistFileUrl ]]; then
  script -f -c "./download_passlist.sh $passlistFileUrl $defaultPasslistFile \
    && python3 ./btcrecover/btcrecover.py --dsw --data-extract-string $walletString \
    --passwordlist $defaultPasslistFile --enable-gpu" \
    $pipe
elif [[ -n $passlistFile ]]; then
  script -f -c "python3 ./btcrecover/btcrecover.py --dsw --data-extract-string $walletString \
    --passwordlist $passlistFile --enable-gpu" \
    $pipe
else
  echo "Incorrect parameters"
  exit 1
fi
