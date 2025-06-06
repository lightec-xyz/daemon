## ZKBTC JsonRpc

Zkbtc JSON-RPC API for interaction with zkbtc node,the default json rpc server use 9977 prot.

### zkbtc_version

return the current zkbtc node info

*example*

      curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"'"zkbtc_version"'","params":[],"id":1}' http://127.0.0.1:9977

### zkbtc_txesByAddr

return an address transaction

*Parameters*

    // query an address all redeem transactions
    "params":["0x5eed85149D7C3d74d28C9b164b210a20e749199c","redeem"]
    // query an address all deposit transactions
    "params":["0x5eed85149D7C3d74d28C9b164b210a20e749199c","deposit"]

*example*

       curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"'"zkbtc_txesByAddr"'","params":["0x5eed85149D7C3d74d28C9b164b210a20e749199c","redeem"],"id":1}' http://127.0.0.1:9977

### zkbtc_transaction

query transaction by giving a tx hash

*example*

    curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"'"zkbtc_transaction"'","params":["0x6deff065bbaf2c9e9c12faf1d841d1f0b96502a20e6e5a864cc398cf6d54d6e4"],"id":1}' http://127.0.0.1:9977

### zkbtc_proofInfo

return zk proof by giving a tx hash

*example*

    curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"'"zkbtc_proofInfo"'","params":["0xdc66dae7e4e27a61884791706377561c60fafa5b10160d970dbff0ebd552657d"],"id":1}' http://127.0.0.1:9977


    