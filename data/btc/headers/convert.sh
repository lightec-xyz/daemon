#!/bin/sh

set -e

convert_time() {

  if [ $(uname) = 'Darwin' ]; then
      d=$(TZ=UTC date -jf "%Y-%m-%d%H:%M:%S" $1 "+%s")
  elif [ $(uname) = 'Linux' ]; then
      echo "not implemented, abort"
      # should be sth. like date -d $1 ...
      exit 1
  fi

  timestamp=$(printf "%x\n" $d | ./endian.sh)
}

convert() {

  content_lines=$(cat $1)
  last_hash=$3

  IFS=$'\n'
  for content in $content_lines; do

    height=$(echo $content | awk '{print $1}')

    version_decimal=$(echo $content | awk '{print $2}')
    version=$(printf "%08x\n" $version_decimal | ./endian.sh)

    this_hash=$(echo $content | awk '{print $3}' | ./endian.sh)
    merkle_root=$(echo $content | awk '{print $4}' | ./endian.sh)
    
    timestamp_str=$(echo $content | awk '{print $5}')
    convert_time $timestamp_str
    
    bits_decimal=$(echo $content | awk '{print $6}')
    bits=$(printf "%08x\n" $bits_decimal | ./endian.sh)

    nonce_decimal=$(echo $content | awk '{print $7}')
    nonce=$(printf "%08x\n" $nonce_decimal | ./endian.sh)

    header=$version$last_hash$merkle_root$timestamp$bits$nonce

    echo $header >> $2
    last_hash=$this_hash
  
   done
}

convert $1 $2 $3

cat $2 >> headers.txt
