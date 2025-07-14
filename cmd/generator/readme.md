# Testnet Generator

    mkdir testnet && cd testnet

download sources code

    git clone https://github.com/lightec-xyz/RLPark.git 

    git clone https://github.com/lightec-xyz/btc_provers.git 

    git clone https://github.com/lightec-xyz/provers.git

    git clone https://github.com/lightec-xyz/gMPTark.git 

    git clone https://github.com/lightec-xyz/daemon.git

    cd daemon
    git checout -b testnet origin/mainnet


setup

    cd cmd/setup
    go build -tags="prod"
    ./setup  --datadir /opt/circuitsetup/testnet --srsdir /opt/aztecSrs/aztec_bn254 run --group all --chainId 11155111 --zkbtcBridgeAddr 0x95E728c449ffF859CfbC100835e53064E6134f33 --icpPublickey 02971351ad0a4e80b4d61003a152c746bde6d7ac5cba52466727c611fdc8c20f5b 

generator

    cd cmd/generator
    go build -tags="prod"
    
    ./generator --config ./config.json run

 config
 
    {
        "url": "http://3.137.171.200:10811",
        "maxNums": 2,
        "cacheCap":0,
        "network": "testnet",
        "disableVerifyZkFile":true,
        "ethSetupDir": "/opt/circuitsetup/testnet",
        "btcSetupDir": "/opt/circuitsetup/testnet",
        "datadir": "/opt/lightec/testnet/.generator",
        "mode": "client"
    }


