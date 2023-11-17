#!/bin/sh

set -e
mkdir -p tmp
previous_hash=0000000000000000000000000000000000000000000000000000000000000000

# split the files
line_count=$(cat -n sorted.txt | tail -1 | awk '{print $1}')
let "step = $line_count / 20"
echo "line count is $line_count, step size is $step"

# not sure why but have to run below scripts manually
for i in {0..20}
do
  let "head = $step * $i"
  let "tail = $line_count - $head"

  echo "tail -$tail sorted.txt | head -$step > tmp/sorted.txt.$i"
  tail -$tail sorted.txt | head -$step > tmp/sorted.txt.$i

  echo "./convert.sh tmp/sorted.txt.$i tmp/headers.txt.$i $previous_hash &"
  ./convert.sh tmp/sorted.txt.$i tmp/headers.txt.$i $previous_hash &

  previous_hash=$(tail -1 tmp/sorted.txt.$i | awk '{print $3}' | ./endian.sh)
done

