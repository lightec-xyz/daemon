# daemon

A node for a cross-chain bridge between Ethereum and Bitcoin implemented in the Lightning protocol

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

### Command

daemon

    // run daemon process
    ./daemon run 

    // add remote worker to daemon 
    ./daemon addWorker http://127.0.0.1:8485 1

    // stop daemon
    ./daemon stop 

proof worker

    ./proof  run

    
    
    


    