#!/bin/sh
# Download all bitcoin blocks from blockchair.com

set -e

mkdir -p blocks
cd blocks

download() {
  echo "downloading ${1} from https://gz.blockchair.com"
  curl https://gz.blockchair.com/bitcoin/blocks/blockchair_bitcoin_blocks_${1}.tsv.gz > blockchair_bitcoin_blocks_${1}.tsv.gz
  sleep 1
}

start_date="20250507"

for date in $(cat ../dates.txt); do

  if [[ $date < $start_date ]]; then
    continue;
  fi

  if ! [ -s blockchair_bitcoin_blocks_${date}.tsv.gz ]; then
    if [ -f blockchair_bitcoin_blocks_${date}.tsv.gz ]; then
      echo "found empty file blockchair_bitcoin_blocks_${date}.tsv.gz, removed"
      rm blockchair_bitcoin_blocks_${date}.tsv.gz
      echo "re-downloading"
    fi
    download $date
  fi
done