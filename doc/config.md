Here is a template for the daemon node configuration file. You can change it according to your needs.

    {
        "datadir": "/opt/lightec/devnet/.daemon",
        "btcMainnetPath":"/opt/circuit/headers.txt",
        "network": "testnet",
        "rpcbind": "0.0.0.0",
        "rpcport": "10977",
        "wsPort": "8880",
        "btcUser": "lightec",
        "sgxUrl": "http://127.0.0.1:8977",
        "btcPwd": "lighteciwueh7297sd",
        "btcUrl": "http://127.0.0.1:9935",
        "ethUrl": "http://127.0.0.1:9899",
        "beaconUrl": "http://127.0.0.1:8365",
        "oasisUrl": "https://testnet.sapphire.oasis.dev",
        "discordHookUrl": "",
        "ethPrivateKey": "",
        "btcInitHeight":79775,
        "btcGenesisHeight":74592 ,
        "icpWalletAddress":"xbsjj-qaaaa-aaaai-aqamq-cai",
        "beaconInitSlot": 4161092,
        "genesisBeaconSlot": 4128768,
        "ethInitHeight":3749086,
        "enableLocalWorker": false,
        "btcReScan": false,
        "disableBtcAgent":false,
        "disableEthAgent":false,
        "disableBeaconAgent":false,
        "disableFetch": false,
        "disableLipP2p": true,
        "p2pBootstraps": [
            ""
        ]
    }

Generator config template

    {
        "url": "http://127.0.0.1:10811",
        "maxNums": 2,
        "cacheCap": 1,
        "network": "testnet",
        "datadir": "/opt/lightec/testnet/.generator",
        "ethSetupDir": "/opt/circuitsetup/testnet/",
        "btcSetupDir": "/opt/circuitsetup/testnet/",
        "disableVerifyZkFile":false,
        "model": "client"
    }
