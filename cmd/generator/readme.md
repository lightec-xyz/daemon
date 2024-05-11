# Generator

## Build

    git clone https://<Token>@github.com/lightec-xyz/daemon
    git checkout -b red_dev orign/red_dv
    cd  daemon/cmd/generator && go build


## Environment
1. download circuit config data


    
    // maybe you should Compressed file before download
    scp -r red@58.41.9.129:/opt/lightec/circuit_data/beta1/circuits <lcoal circuit param file dir>

2. set ZkParameterDir environment variables



    export ZkParameterDir =<local circuit param file dir>

3. generator run config file



    {
        "url": "https://test.apps.zkbtc.money/api",
        "maxNums": 1,
        "network": "testnet",
        "datadir": "<local genetator storagte data dir>",
        "model": "client"
    }


## Run

    ./generator --config ./testnet.json run

