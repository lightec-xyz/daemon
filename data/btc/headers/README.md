## Jump start the dates.txt file if not already there

Using this command to generate date list:

```
cat file_list.txt | sed -e 's/^.*blocks_//' | sed -e 's/\.tsv\.gz.*$//' > dates.txt
```

Initial file_list.txt is from [blockchair](https://gz.blockchair.com/bitcoin/blocks). Note that the dates are _not_
continuous initially.

## Update instructions

Since the initial data collection has been completed, below listed are instructions for incremental updates.

- Check the value of `start_date` before running each of `udpate_dates.sh`, `download_blocks.sh`, `sort_blocks.sh`.
- Use the `update_dates.sh` script to update the `dates.txt` file.
- Unzip the `headers.txt.gz` file with the `gunzip` command.
- Use the `download_blocks.sh` script to download from blockchair.com the bitcoin blocks for each new date listed in the
  `dates.txt` file.
- Use the `sort_blocks.sh` script to extract needed fields from downloaded files and then sort them according to block
  heights. Results are saved in the `sorted.txt` file.
- Retrieve the hash of last block before this update with this command
  ` tail -1 headers.txt | head -c 160| xxd -r -p | openssl sha256 -binary | openssl sha256 -hex `.
- Use the `convert.sh sorted.txt new_headers.txt $last_hash` to convert the entries. It then appends the newly converted
  entries to the `headers.txt` file.
- Run `go test` to quickly verify if all the entries in the file `headers.txt` links correctly. It also tries to fix any
  missing header or wrong headers. Fixed file is save in `headers.txt.fixed`. You may rename it to `headers.txt` and run
  `go test` again. If no error is reported, the file is good to use.

`splitting` and `merging` are for initial data collecting only, executed between `sorting` and optional `verification`.

- Use the `split.sh` script to split the `sorted.txt` files into around 20 smaller files, and convert the entries in
  each new file into block headers by calling the `convert.sh` script. This script might run for **a few hours**
  although the conversion runs in parallel in the background. You need to wait for all background processes to finish.
  The `convert.sh` script also verifies all the generated entries and tries to fix if any errors. If any errors are
  reported and/or fixed, you need to **handle them** before proceeding to the `merge.sh` script.
- Use the `merge.sh` script to merge the smaller header files into the `headers.txt` file.

### Known issues

- There are some errors in the downloaded data. `go test` may fix them.

## Programtic retrieval

- Use the `getheaders` p2p message to retrieve 2000 headers at once.
