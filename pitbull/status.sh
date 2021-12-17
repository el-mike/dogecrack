#!/bin/bash

SUCCESS_CODE=0
SUCCESS_PREFIX="SUCCESS: "

# Arbitrary code values for potentially "conflict-free" codes:
RUNNING_CODE=50
RUNNING_PREFIX="RUNNING: "

INTERRUPTED_CODE=51
INTERRUPTED_PREFIX="STOPPED: "

FINISHED_CODE=52
FINISHED_PREFIX="FINISHED: "

lastLine=$(tail -1 ./out_btcrecover.txt)

# If btcrecover succeeded, last line contains: 'Password found: password'
successCheck=$(cat ./out_btcrecover.txt | grep 'found')

if [[ $successCheck ]]
then
  echo $SUCCESS_PREFIX $lastLine

  exit $SUCCESS_CODE
fi

interruptionCheck=$(cat ./out_btcrecover.txt | grep 'Interrupted')

if [[ $interruptionCheck ]]
then
  # We want to echo line with the interruption message, as it contains password line count.
  # It can be used later to start btcrecover on given line.
  echo $INTERRUPTED_PREFIX $interruptionCheck

  exit $INTERRUPTED_CODE
fi

# If out file does not contain "Password found" or "Interrupted" line, check if still going.
# If so, return RUNNING_PREFIX with ps info.
runningCheck=$(ps -A | grep python)

if [[ $runningCheck ]]
then
  echo $RUNNING_PREFIX $runningCheck

  exit $RUNNING_CODE
fi

# Otherwise, it finished, and did not find a password.

echo $FINISHED_PREFIX $lastLine
exit $FINISHED_CODE
