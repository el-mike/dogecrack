#!/bin/bash

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

WAITING_STATUS="WAITING"
RUNNING_STATUS="RUNNING"
FINISHED_STATUS="FINISHED"
SUCCESS_STATUS="SUCCESS"

# If btcrecover succeeded, last line contains: 'Password found: password'.
# We grep the entire file though, to be sure that it works in case of some
# unnecessary ouput is added to the end of the file.
# We add stderr redirect at the end of "cat" to avoid printing errors - we only care about
# exit code checked below (for WAITING status).
successCheck=$(cat ./progress_view.txt 2>/dev/null | grep '[P]assword found')

# If successCheck returned exit code 1 (meaning progress_view.txt could not be read,
# i.e. does not exist), it means that pitbull has not been run yet - therefore we return
# WAITING status.
if [[ $? -eq 1 ]]; then
  echo $WAITING_STATUS
  exit 0
fi

if [[ $successCheck ]]; then
  echo $SUCCESS_STATUS
  exit 0
fi

# If out file does not contain "Password found" line, check if still going.
# If so, return RUNNING status.
# In order to check if pitbull is still running, we get all terminal processes
# with "ps l", and then will search for process that includes "pitbull" command.
#
# Square brackets around the "p" letter excludes grep itself from search results.
runningCheck=$(ps l | grep '[r]un_pitbull')

if [[ $runningCheck ]]; then
  echo $RUNNING_STATUS
  exit 0
fi

# Otherwise, it finished, and did not find a password.
echo $FINISHED_STATUS
exit 0
