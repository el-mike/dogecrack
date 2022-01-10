#!/bin/bash

ord() {
  LC_CTYPE=C printf '\n%d' "'$1"
}

# Checks if given line is a progress bar from gdown tool.
# Regex tests for progress bar with "| |" at the start and end,
# and for "B/s" download speed suffix.
is_gdown_progress_bar_line() {
  local regex='.*\|.*\|.*B/s.*'

  if [[ $1 =~ $regex ]]; then
    echo 1
  else
    echo 0
  fi
}

# Checks if given line is ETA progress line.
# Regex tests for ETA prefix ands "couting" word that appears at the end of the line.
is_counting_line() {
  local regex='.*ETA.*counting.*'

  if [[ $1 =~ $regex ]]; then
    echo 1
  else
    echo 0
  fi
}

# Checks if given line is a progress bar line.
# Regex tests for the progress bar widget (square brackets with someting inside)
# and ETA information suffix.
is_progress_bar_line() {
  local regex='\[.*\].*ETA.*'

  if [[ $1 =~ $regex ]]; then
    echo 1
  else
    echo 0
  fi
}
