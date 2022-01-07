#!/bin/bash

view_file='progress_view.txt'

# clear file - we want to create new progress view every time we run the reader.
> $view_file

d=$'\r'

ord() {
  LC_CTYPE=C printf '\n%d' "'$1" >> /dev/pts/3
}

while read -r -d $'\r' line
do
  if [[ $line =~ .*counting.* ]]
  then
    sed -i '$ d' $view_file
  fi

  # if [[ $line =~ $'\r' ]]
  # then
  #   echo "FOUND!" >> /dev/pts/3
  # fi

  # ord $line

  # echo "$line" >> /dev/pts/3
  echo "$line" >> $view_file
done <btcrecover_out
 