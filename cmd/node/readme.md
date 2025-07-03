## Command

run

    ./node --config ./node.json run

jwt

    ./node --config ./local.json jwt

scRoot

    ./node --config devnet.json  scRoot --period 939 --genesisSlot 7692288

readVk

    ./node proof readVk --file /opt/lightec/opt2/circuitsetup/devnet/redeem.vk

export proof

    //./node proof export --proof /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/blockchain_0_901151.proof --witness /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/blockchain_0_901151.wtns --name btcduperrecursive_899135_901151  --datadir <path>

    ./node proof import --proof <path> --witness <path> --name <name> --datadir <path>
nonce

        ./node --config ./devnet.json miner nonce --addr 0xb4183bB52E44C6861AEF3B626eb2195288AfCa2f --nonce 1239

filestore

    ./node --config ./devnet.json proof filestore --height 8000




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
        
    




