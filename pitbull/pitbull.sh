#!/bin/bash

dirname=$(dirname "$0")

passlistFileName='passlist.txt'
outFileName='out_btcrecover.txt'

pipe='btcrecover_out'

download_google_drive() {
  gdown  https://drive.google.com/uc?id=$1 -O $passlistFileName
}

download() {
  wget $1
}

detachedMode=0

while getopts f:w:g:d flag
do
    case "${flag}" in
        f) fileUrl=${OPTARG};;
        g) googleFileId=${OPTARG};;
        w) walletString=${OPTARG};;
        d) detachedMode=1;;
    esac
done

echo "Wallet string: $walletString"

# Input args validation
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

# Output pipe setup
mkfifo "$pipe"

if [[ $googleFileId ]]
then
  echo "GoogleDrive file source - using gdown..."
  download_google_drive $googleFileId
else 
  echo "Using wget..."
  download $fileUrl
fi

if [[ $detachedMode -eq 1 ]]
then
  echo "Running in detached mode..."
  
  python3 $dirname/btcrecover/btcrecover.py --dsw \
    --data-extract-string $walletString \
    --passwordlist $passlistFileName \
    --enable-gpu \
    &> $outFileName & # runs the process 
else
  python3 $dirname/btcrecover/btcrecover.py --dsw \
    --data-extract-string $walletString \
    --passwordlist $passlistFileName \
    --enable-gpu \
    2>&1 | tee $outFileName
fi
