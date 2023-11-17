# daemon

A node for a cross-chain bridge between Ethereum and Bitcoin implemented in the Lightning protocol,you can learn more
here [zkbtc](https://www.zkbtc.money/)

## Prerequisites

To gather blockchain data, the following blockchain node JSON RPC methods must be supported

1. bitcoin full node [Bitcoin](https://github.com/bitcoin/bitcoin),we recommend using node version 28.1
2. ethereum beacon chain full node [Prysm](https://github.com/prysmaticlabs/prysm),we recommend using node version
   v6.0.0
3. ethereum execution chain full node [Ethereum](https://github.com/ethereum/go-ethereum),we recommend using node
   version
   v1.15.10

## Hardware requirements

1. Daemon just gather data, ordinary servers are sufficient, there are no special requirements
2. The generator requires significant CPU and memory resources to create zero-knowledge proofs, along with additional
   disk space to store the zero-knowledge parameter file ,the following are the machine parameters we tested

Minimum machine requirements,other tested
machines [example](https://github.com/lightec-xyz/provers/blob/e016c3de22f37540cc489b72d479a282ca87439a/README.md?plain=1#L116)

    CPU   : 24 cores
    Memory: 64GB
    Disk  : 250GB

## Setup

1.Build setup program

    git clone https://github.com/lightec-xyz/daemon
    cd  daemon/cmd/setup && go build

2.Download [aztecSrs]() to generate zero-knowledge parameter file

3.Run setup program

    ./setup  --datadir <circuitDir> --srsdir <srsDir> run --group all --chainId 17000 --zkbtcBridgeAddr 0x49793ff075b696e6bef1b85e4e85fe669041b312 --icpPublickey 023f203422be55a3576f46dc6770bdc7865a126381c1963a2d82b49f4158409a2e

## Build

Build zkbtc node daemon ( gather relevant transaction data from the blockchain)

    git clone https://github.com/lightec-xyz/daemon
    cd  daemon/cmd/node && go build

Build zkbtc daemon generator (the program that generates zero-knowledge proofs)

    git clone https://github.com/lightec-xyz/daemon
    cd  daemon/cmd/generator && go build

## Run

The daemon default storage location is: ***~/.daemon***,run zkbtc daemon ,you can find other field detailed
explanations [here](./config.md).

    // config.json
    {
        "btcUrl": "http://127.0.0.1:8332",
        "ethUrl": "http://127.0.0.1:8545",
        "beaconUrl": "http://127.0.0.1:30814",
    }

    // run node
    ./node --config ./config.json run

The generator default storage location is: ***~/.generator*** ,run zkbtc generator,config template [here](./config.md).

    //config.json
    {
        "url": "http://127.0.0.1:8970",
        "maxNums": 1,
        "ethSetupDir": "<circuitSetupPath>",
        "btcSetupDir": "<circuitSetupPath>",
    }

    // run generator
    ./generator --config ./config.json run

    
    
    


    