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
    ./setup  --datadir /opt/circuitsetup/mainnet --srsdir /opt/aztecSrs/aztec_bn254 run --group all --chainId 1 --zkbtcBridgeAddr 0xB86E9A8391d3df83F53D3f39E3b5Fce4D7da405d --icpPublickey 03183007b9afcfa519871885380d4dfd1144269d8050ec2a51992065af2a87d3df 

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


