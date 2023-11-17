#!/bin/bash
# from: https://stackoverflow.com/questions/25701265/how-to-generate-a-list-of-all-dates-in-a-range-using-the-tools-available-in-bash
# with "until"
start_date="2025-05-07"
today=$(TZ=UTC date -I)
until [[ $start_date = $today ]]; do
    echo "adding $start_date"
    echo "$start_date" | sed -e 's/-//g'  >> dates.txt
    if [ $(uname) = 'Darwin' ]; then
        start_date=$(date -j -v+1d -f %Y-%m-%d $start_date +%Y-%m-%d)
    elif [ $(uname) = 'Linux' ]; then
        start_date=$(date -I -d "$start_date + 1 day")
    fi
    # sleep 1
done