#!/bin/bash

# This scripts returns the progress information in a form of "done of to_be_done".
# Progress is taken from btcrecover's progress line output.

viewFile='progress_view.txt'

source ./helpers.sh

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
  echo "Progress not found: $lastLine"
fi
