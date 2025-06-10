## Command

run

    ./node --config ./node.json run

jwt

    ./node --config ./local.json jwt

scRoot

    ./node --config devnet.json  scRoot --period 939 --genesisSlot 7692288

readVk

    ./node readVk --file /opt/lightec/opt2/circuitsetup/devnet/redeem.vk

export proof

    //example ./node exportProof --name btcduperrecursive_84010_84024 --datadir ./daemon --proof btcduperrecursive_84010_84024.proof --witness btcduperrecursive_84010_84024.witness 
    
    ./node exportProof --proof <path> --witness <path> --name <name> --datadir <path>
nonce

    ./node --config ./node.json miner --miner <0x> --nonce 1000



## Config

    {
        "datadir": "/Users/red/lworkspace/lightec/daemon/node/test",
        "network": "local",
        "rpcbind": "127.0.0.1",
        "rpcport": "9870",
        "btcUrl": "",
        "ethUrl": "https://ethereum-holesky-rpc.publicnode.com",
        "beaconUrl": "",
        "ethPrivateKey": "",
        "btcInitHeight": 	2585500,
        "ethInitHeight": 1298160,
        "enableLocalWorker": false,
        "autoSubmit": true
    }
        
    




