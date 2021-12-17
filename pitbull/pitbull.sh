#!/bin/bash

passFileName='pass.txt'
outFileName='out_btcrecover.txt'

download_google_drive() {
  gdown  https://drive.google.com/uc?id=$1 -O $passFileName
}

download() {
  wget $1
}

while getopts f:w:g: flag
do
    case "${flag}" in
        f) fileUrl=${OPTARG};;
        g) googleFileId=${OPTARG};;
        w) walletString=${OPTARG};;
    esac
done

echo "Wallet string: $walletString"

if [[ -z $fileUrl && -z $googleFileId ]]
then
  echo "Passlist source missing"
  exit 1
fi

if [[ -z $walletString ]]

then
  echo "Wallet string missing"
  exit 1
fi

if [[ $googleFileId ]]
then
  echo "GoogleDrive file source - using gdown..."
  download_google_drive $googleFileId
else 
  echo "Using wget..."
  download $fileUrl
fi


python3 /app/btcrecover/btcrecover.py --dsw \
  --data-extract-string $walletString \
  --passwordlist ./pass.txt \
  --enable-gpu \
  &> $outFileName & # runs the process 


