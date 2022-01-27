#!/bin/bash

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

source ./variables.sh

WAITING_STATUS="WAITING"
RUNNING_STATUS="RUNNING"
FINISHED_STATUS="FINISHED"
SUCCESS_STATUS="SUCCESS"

# If progress_view.txt does not exist, it means that pitbull has not been run yet,
# therefore we return WAITING status.
if [ ! -f "./$viewFile" ]; then
  echo $WAITING_STATUS
  exit 0
fi

# If btcrecover succeeded, last line of progress_view.txt contains:
# 'Password found: <password>'. We check the entire file though, to be sure
# that it works in case of some unnecessary ouput is added to the end of the file.
# We add stderr redirect to avoid printing errors and flooding the output.
successCheck=$(cat ./$viewFile 2>/dev/null | grep '[P]assword found')

if [[ $successCheck ]]; then
  echo $SUCCESS_STATUS
  exit 0
fi

# If out file does not contain "Password found" line, check if still going.
# If so, return RUNNING status.
# In order to check if pitbull is still running, we get all terminal processes
# with "ps l", and then will search for process that includes "run_pitbull" command.
#
# Square brackets around the "p" letter excludes grep itself from search results.
# 
# Note: we cannot use 'pitbull' itself as search value, because "status" command
# is also run via main pitbull.sh file, therefore it would always return true. 
runningCheck=$(ps l | grep '[r]un_pitbull')

if [[ $runningCheck ]]; then
  echo $RUNNING_STATUS
  exit 0
fi

# Otherwise, it finished, and did not find a password.
echo $FINISHED_STATUS
exit 0
