## Command

    ./proof --config ./client_config.json run

## Config

    {
        "url": "https://test.apps.zkbtc.money",  // daemon server url
        "maxNums": 1,                            // maximum number of tasks on the machine,default 1
        "network": "local",                      // lightec network name
        "datadir": "/opt/light/.proof,           // generate proof temporary directory and setup directory, should be separated (environment variable settings ？)
        "model": "client"                        // generator mode (client: server-client)
    }

## tips
**please prepare all setup files before starting the proof generator**