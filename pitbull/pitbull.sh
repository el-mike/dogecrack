#!/bin/bash

# This is the main script of Pitbull tool. It's responsible for reading arguments
# and running the scripts in proper mode.
# By default, Pitbull will start new terminal session with tmux, under the name "pitbull".
# 
# Output will be redirected to "progress_view.txt" file via named pipe (btcrecover_out), but you can always
# re-attach terminal session with `tmux a -t "pitbull"` to see the live progress.
# 
# pitbull.sh is an "entry script" - calling any other script directly may no work properly.

# Returns the directory the script exists in, no matter where it was called from.
# It's needed for referencing other scripts and saving files in consistent way.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Proceed to pitbull directory. From now on, all scripts will be referencing the directory
# pitbull.sh exists in.
cd $dirname

source ./variables.sh

runCommand='run'
statusCommand='status'
progressCommand='progress'
outputCommand='output'
killCommand='kill'
errorsCommand='errors'
attachCommand='attach'

# clear file - we want to create new log file view every time we run Pitbull.
> $errLogFile

# command - describes the operation an user agent wants to perform. Available commands
# are listed above.
command=$1

if [[ "$command" == "$statusCommand"  ]]; then
  ./status.sh
  exit $?
elif [[ "$command" == "$progressCommand"  ]]; then
  ./progress.sh
  exit $?
elif [[ "$command" == "$outputCommand"  ]]; then
  ./output.sh
  exit $?
elif [[ "$command" == "$killCommand" ]]; then
  ./kill.sh
  exit $?
elif [[ "$command" == "$errorsCommand" ]]; then
  ./errors.sh
  exit $?
elif [[ "$command" == "$attachCommand" ]]; then
  ./attach.sh
  exit $?
elif [[ "$command" == "$runCommand"  ]]; then
  # Since we are "starting" with positional argument (command), we need to shift
  # it one place to get the optional params properly.
  shift 1

  if [[ $* == *--no-tmux* ]]; then
    ./run_pitbull.sh "$@" 2>> $errLogFile
  else
    tmux new-session -d -s "pitbull" ./run_pitbull.sh "$@" 2>> $errLogFile
  fi
else
  echo "Command: '$command' not recognized"
  exit 1
fi
