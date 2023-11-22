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


refer
* db:https://www.cnblogs.com/orange-CC/p/13212042.html
* jsonRpc2.0:https://github.com/sourcegraph/jsonrpc2
* gin:https://gin-gonic.com/
* cbor:https://github.com/fxamacker/cbor
