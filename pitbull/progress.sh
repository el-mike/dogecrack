#!/bin/bash

# This scripts returns the progress information in a form of "done of to_be_done".
# Progress is taken from btcrecover's progress line output.

# Returns the directory the script exists in, no matter where it was called from.
dirname=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $dirname

source ./helpers.sh

viewFile='progress_view.txt'
errorMessage='NO_PROGRESS_AVAILABLE'

# If progress_view.txt does not exist, we return fallback message.
if [ ! -f "./$viewFile" ]; then
  echo $errorMessage
  exit 0
fi

lastLine=$(tail -1 "$viewFile")

progress_bar_step=$(is_progress_bar_line "$lastLine")

if [[ $progress_bar_step -eq 1 ]]; then
  # Regex searches from the beginning of the line till "[" character,
  # which is the first character of the visual progress bar "widget". 
  regex='^(.*)\['

  if [[ $lastLine =~ $regex ]]; then
    match=${BASH_REMATCH[0]}
    progress=${match::-2}

    echo "$progress"
  fi

else
  echo $errorMessage
fi
