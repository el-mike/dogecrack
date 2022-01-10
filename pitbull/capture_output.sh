#!/bin/bash

# This script reads the named pipe that pitbull writes into, and tries to
# recreate the output.
# As ETA counting, progress bar and couple other output informations is meant for
# the user to see, we need to use some tricks to capture the output that
# is being constantly overwritten by \r (carriage return) character.

pipe='btcrecover_out'
viewFile='progress_view.txt'

source ./helpers.sh

# clear file - we want to create new progress view every time we run the reader.
> $viewFile

# Since most of the output uses carriage return to overwrite current line,
# we need to use \r delimiter to capture progress properly.
delimiter=$'\r'

while read -r -d $delimiter line
do
  gdown_progress_bar_step=$(is_gdown_progress_bar_line "$line")
  eta_counting_step=$(is_counting_line "$line")
  progress_bar_step=$(is_progress_bar_line "$line")

  # If one of the overwritable line is being outputted, we want to
  # rmeove the last line of view file.
  # Since btcrecover can run for hours, printing every progress flush would
  # make the view file very big. Additionaly, it keeps the progress view clean and readable.
  if [[ $gdown_progress_bar_step -eq 1 || $eta_counting_step -eq 1 || $progress_bar_step -eq 1 ]]; then
    # Remove last line of the file.
    sed -i '$ d' $viewFile
  fi

  echo "$line" >> $viewFile
done < $pipe
 