# Generator

## Build

    makdir <projectDir>

    cd <projectDir>

    git clone https://<Token>@github.com/lightec-xyz/provers.git
    git checkout v0.2.0   

    git cloen https://<Token>@github.com/lightec-xyz/btc_provers.git
    git checkout v0.2.0  

    git clone https://<Token>@github.com/lightec-xyz/reLight.git
    git checkout v0.2.0  


    git clone https://<Token>@github.com/lightec-xyz/daemon
    git checkout -b red_dev orign/red_dv
    cd  daemon/cmd/generator && go build

## Environment

download circuit config data

    // maybe you should Compressed file before download
    scp -r <userName>@58.41.9.129:/opt/lightec/circuit_data/beta1/circuits <lcoal circuit param file dir>

set ZkParameterDir environment variables

    export ZkParameterDir =<local circuit param file dir>

generator run config file

    {
        "url": "https://test.apps.zkbtc.money/api",
        "maxNums": 1,
        "network": "testnet",
        "datadir": "<local genetator storagte data dir>",
        "model": "client"
    }

## Run

    ./generator --config ./testnet.json run

