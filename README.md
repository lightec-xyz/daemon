# daemon

A node for a cross-chain bridge between Ethereum and Bitcoin implemented in the Lightning protocol

***Note:*** The project is continuously in development, and there may be incompatible changes in the code and API.

## Network

Bitcoin Testnet Network:

* OperatorAddress: tb1qalv7aduqdpz9wc4fut3nt44tsf42anleed76yj3el3rgd4rgldvq8yyuhr

Ethereum Holesky Network

* zkBridgeAddress: 0xbdfb7b89e9c77fe647ac1628416773c143ca4b51

## Build

    git clone https://<token>@github.com/lightec-xyz/daemon
    cd  daemon/cmd/daemon && go build

## Run

### Introduce

1. the daemon default storage location is: ***~/.daemon***
2. the daemon config file location is: ***~/.daemon/daemon.json***

When you deploy for the first time, you need to modify the relevant parameters. The following are the required
parameters. You can find other detailed explanations [here](./doc/config.md).

    {
        "btcUrl": "http://127.0.0.1:8332",      / /Bitcoin Core jsonrpc endpoint                           
        "btcUser": "lightec",                   // The bitcoin jsonrpc authentication username, if you have set one
        "btcPwd": "abcd1234",                   // The bitcoin jsonrpc authentication password
        "ethUrl": "http://localhost:8545"       // Ethereum jsonrpc endpoint
    }

### Command

daemon

    // run daemon process
    ./daemon run 

    // add remote worker to daemon 
    ./daemon --rpcbind 127.0.0.1 --rpcport 9780 addWorker ws://127.0.0.1:30001 1

    // stop daemon
    ./daemon stop 

proof worker

    ./proof  run

    
    
    


    