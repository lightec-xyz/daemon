## Daemon 
* store db
  * interface
  * leveldb or rocksdb
  * Serialization and Deserialization
    * borsh，protobuf，cbor，rlp,scale
  * table
* api server
  * interface
  * jsonRpc2.0
  * gin or other
  * client
* scan block
  * bitcoin
  * ethereum
* transaction
  * bitcoin
    * definity,oss
  * ethereum
    * 
* command line
* main struct 


## Test
* bitcoin
  * mainnet:https://bitcoin-mainnet-archive.allthatnode.com
  * testnet:https://bitcoin-testnet-archive.allthatnode.com

      
    curl https://bitcoin-testnet-archive.allthatnode.com \
    --request POST \
    --header 'content-type: text/plain;' \
    --data '{"jsonrpc": "1.0", "id": "curltest", "method": "getblockcount", "params": []}'
    
