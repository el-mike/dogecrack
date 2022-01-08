#!/bin/bash

# This script downloads given passlist file (into passlist.txt) and runs btcrevcover with it.
# It's the main functionality of Pitbull tool.
# It will run as a foreground process, and output progress to TTY. Some additional logs
# may be redirected to stderr (warnings, errors).

dirname=$(dirname "$0")

passlistFileName='passlist.txt'
pipe='btcrecover_out'

# Checks if given file URL is hosted on Google Drive.
is_google_drive_file() {
  local regex='.*drive.google.com.*'

  if [[ $1 =~ $regex ]]; then
    echo 1
  else
    echo 0
  fi
}

# Downloads a file from Google Drive using gdown.
# We are using custom tools, as Google Drive performs additional checks while
# accessing and downloading files.
download_from_google_drive() {
  gdown  $1 -O $passlistFileName
}

# Downloads a file using wget.
download() {
  wget $1
}


while getopts f:w: flag
do
    case "${flag}" in
        f) fileUrl=${OPTARG};;
        w) walletString=${OPTARG};;
    esac
done

echo "Wallet string: $walletString"

# Input args validation.
if [[ -z $fileUrl ]]; then
  echo "Passlist source missing"
  exit 1
fi

if [[ -z $walletString ]]; then
  echo "Wallet string missing"
  exit 1
fi

isGoogleDriveFile=$(is_google_drive_file "$fileUrl")

# Downloadind passlist file.
if [[ $isGoogleDriveFile -eq 1 ]]; then
  echo "GoogleDrive file source - using gdown..."
  download_from_google_drive $fileUrl
else 
  echo "Using wget..."
  download $fileUrl
fi

# Output capture setup.
# If pipe exists, remove it - it ensures that no other agent is
# reading from the output pipe.
if [ -p $pipe ]; then
  rm "$pipe"
fi

mkfifo "$pipe" && ./output_reader.sh &

script -f -c "python3 $dirname/resources/btcrecover/btcrecover.py --dsw \
  --data-extract-string $walletString \
  --passwordlist $passlistFileName \
  --enable-gpu" \
  $pipe
