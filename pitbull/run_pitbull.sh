#!/bin/bash

# This script runs btcrecover with appropriate arguments, decided by the options passed by the user.
# Based on the arguments provided, it will be run in one of three modes:
# - with "-t" flag, it will use provided tokenlist file and pass it to btcrecover
# - with "-f" flag, it will use provided passlist file and pass it to btcrecover
# - with "u" flag, it will download passlist file from given URL and pass downloaded file to btcrecover
# It will run as a foreground process, and output progress to TTY. Some additional logs
# may be redirected to stderr (warnings, errors).
# By having the actual run call in separated file, we can easily modify the way
# Pitbull will be run (for example in detached TTY created with tmux).

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

source ./variables.sh

while [[ $# -gt 0 ]]; do
  case $1 in
    -t)
      tokenList=$2
      # We need to shift twice for both arg name and value
      shift
      shift;; # end of case
    -f)
      passlistFile=$2
      shift
      shift;;
    -u)
      passlistFileUrl=$2
      shift
      shift;;
    -w)
      walletString=$2
      shift
      shift;;
    --skip)
      skipCount=$2
      shift
      shift;;
    --length-min)
      lengthMin=$2
      shift
      shift;;
    --length-max)
      lengthMax=$2
      shift
      shift;;
    # --no-tmux is handled in the main script, but it's passed down here - it can be safely ignored.
    # it also doesn't take a value, so we only shift once.
    --no-tmux)
      shift;;
    -*|--*)
      echo "Unknown option: $1"
      exit 1
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

skipArg=''
lengthMinArg=''
lengthMaxArg=''

# If not empty, build an argument.
# Note that dynamic arguments include a single space at the beginning, as to not add unnecessarily padding when absent.
if [[ ! -z $skipCount ]]; then
  skipArg=" --skip $skipCount"
fi

if [[ ! -z $lengthMin ]]; then
  lengthMinArg=" --length-min $lengthMin"
fi

if [[ ! -z $lengthMax ]]; then
  lengthMaxArg=" --length-max $lengthMax"
fi

dynamicArgs="$skipArg$lengthMinArg$lengthMaxArg"

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

if [[ -n $tokenList ]]; then
  # We need to use -e flag to preserve newline characters from tokenlist argument.
  echo -e "$tokenList" > "$defaultTokenlistFile"

  script -f -c "python3 ./btcrecover/btcrecover.py --dsw --data-extract-string $walletString \
    --tokenlist $defaultTokenlistFile --enable-gpu$dynamicArgs" \
    $pipe
elif [[ -n $passlistFileUrl ]]; then
  script -f -c "./download_passlist.sh '$passlistFileUrl' '$defaultPasslistFile' \
    && python3 ./btcrecover/btcrecover.py --dsw --data-extract-string $walletString \
    --passwordlist $defaultPasslistFile --enable-gpu$dynamicArgs" \
    $pipe
elif [[ -n $passlistFile ]]; then
  script -f -c "python3 ./btcrecover/btcrecover.py --dsw --data-extract-string $walletString \
    --passwordlist $passlistFile --enable-gpu$dynamicArgs" \
    $pipe
else
  echo "Incorrect parameters"
  exit 1
fi
