#!/bin/bash

# This script is responsible for downloading passlist file and saving it for later use.

fileUrl=$1
passlistFile=$2

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
# It saves file in current working directory (that should be set in pitbull.sh).
download_from_google_drive() {
  gdown  $1 -O $passlistFile
}

# Downloads a file using wget.
# It saves file in current working directory (that should be set in pitbull.sh).
download() {
  wget $1
}

isGoogleDriveFile=$(is_google_drive_file "$fileUrl")

# Downloadind passlist file. We use various tools for downloading depending on
# the file storage provider.
if [[ $isGoogleDriveFile -eq 1 ]]; then
  echo "GoogleDrive file source - using gdown..."
  download_from_google_drive $fileUrl
else 
  echo "Using wget..."
  download $fileUrl
fi
