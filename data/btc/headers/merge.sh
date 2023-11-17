#!/bin/sh

set -e

merge() {
  for i in {0..20}
  do
    echo "cat tmp/headers.txt.$i >> headers.txt"
    cat tmp/headers.txt.$i >> headers.txt
  done
  echo "done merge"
}

rm -f headers.txt
merge
