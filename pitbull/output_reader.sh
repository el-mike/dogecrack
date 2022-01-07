#!/bin/bash

# This scripts reads the named pipe that pitbull writes into, and tries to
# recreate the  output.
# As ETA counting, progress bar and couple other output informations is meant for
# the user to see, we need to use some tricks to capture the output that
# is being constantly overwritten by \r (carriage return) character.

pipe='btcrecover_out'
viewFile='progress_view.txt'

# clear file - we want to create new progress view every time we run the reader.
> $viewFile

ord() {
  LC_CTYPE=C printf '\n%d' "'$1"
}

# Checks if given line is ETA progress line.
# Regex tests for ETA prefix ands "couting" word that appears at the end of the line.
is_counting_line() {
  local regex='.*ETA.*counting.*'

  if [[ $1 =~ $regex ]]
  then
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

  if [[ $1 =~ $regex ]]
  then
    echo 1
  else
    echo 0
  fi
}

# Since most of the output uses carriage return to overwrite current line,
# we need to use \r delimiter to capture progress properly.
delimiter=$'\r'

while read -r -d $delimiter line
do
  eta_counting_step=$(is_counting_line "$line")
  progress_bar_step=$(is_progress_bar_line "$line")

  # If one of the overwritable line is being outputted, we want to
  # rmeove the last line of view file.
  # Since btcrecover can run for hours, printing every progress flush would
  # make the view file very big. Additionaly, it keeps the progress view clean and readable.
  if [[ $eta_counting_step -eq 1 || $progress_bar_step -eq 1 ]]
  then
    # Remove last line of the file.
    sed -i '$ d' $viewFile
  fi

  echo "$line" >> $viewFile
done < $pipe
 