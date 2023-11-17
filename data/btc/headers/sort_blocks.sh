#!/bin/sh

set -e

# make it restartable
start_date=20250507

tmp=tmp.dump
rm -f $tmp

for date in $(cat dates.txt); do
  if [[ $date -lt $start_date ]]; then
    continue
  fi

  if ! [ -s blocks/blockchair_bitcoin_blocks_${date}.tsv.gz ]; then
    echo "file does not exist or empty: blockchair_bitcoin_blocks_${date}.tsv.gz. Abort"
    rm -f $recordfile
    exit 1
  fi

  gunzip --stdout --keep blocks/blockchair_bitcoin_blocks_${date}.tsv.gz | grep  -e '^id' -v | awk '{print $1 " " $10 " " $2 " " $13 " " $3 $4 " " $15 " " $14 }' >> $tmp
done

sort -k1 -u -n $tmp > sorted.txt
rm $tmp
