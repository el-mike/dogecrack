#!/bin/bash

passFileName='pass.txt'

download_google_drive() {
  gdown $1 -O $passFileName
}

download() {
  wget $1
}

while getopts f:w:g flag
do
    case "${flag}" in
        f) fileUrl=${OPTARG};;
        w) walletString=${OPTARG};;
        g) usingGoogleDrive=1;;
    esac
done

echo "File URL: $fileUrl"
echo "Wallet string: $walletString"

if [[ -z $fileUrl ]]
then
  echo "File URL missing!"
  exit 1
fi

if [[ -z $walletString ]]

then
  echo "Wallet string missing!"
  exit 1
fi

if [[ $usingGoogleDrive -eq 1 ]]
then
  echo "Using Google Drive"
  download_google_drive $fileUrl
else 
  download $fileUrl
fi


python3 /app/btcrecover/btcrecover.py --dsw --data-extract-string $walletString --passwordlist ./pass.txt --enable-gpu
