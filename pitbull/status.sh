#!/bin/bash

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

SUCCESS_CODE=0
SUCCESS_STATUS="SUCCESS"

# Arbitrary code values for potentially "conflict-free" codes:
RUNNING_CODE=50
RUNNING_STATUS="RUNNING"

FINISHED_CODE=51
FINISHED_STATUS="FINISHED"

# If btcrecover succeeded, last line contains: 'Password found: password'.
# We grep the entire file though, to be sure that it works in case of some
# unnecessary ouput is added to the end of the file.
successCheck=$(cat ./progress_view.txt | grep 'Password found')

if [[ $successCheck ]]
then
  echo $SUCCESS_STATUS
  exit $SUCCESS_CODE
fi

# If out file does not contain "Password found" line, check if still going.
# If so, return RUNNING status.
# In order to check if pitbull is still running, we get all terminal processes
# with "ps l", and then will search for process that includes "pitbull" command.
#
# Square brackets around the "p" letter excludes grep itself from search results.
runningCheck=$(ps l | grep '[p]itbull')

if [[ $runningCheck ]]
then
  echo $RUNNING_STATUS
  exit $RUNNING_CODE
fi

# Otherwise, it finished, and did not find a password.
echo $FINISHED_STATUS
exit $FINISHED_CODE
