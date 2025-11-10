# circuit Setup

#### Command

generate single circuit setup:(beaconSyncCommittee,ethTxInEth2,ethBeaconHeader,
ethFinalityHeader,ethRedeem,btcBlockChain,btcBlockDepth,btcTxInChain,btcDeposit,btcRedeem)

    ./setup  --datadir ./testdir --srsdir ./srsdir run --type ethBeaconHeader 

batch generate circuit setup:(all, bitcoin, beacon, ethereum,txes)

    ./setup  --datadir ./testdir --srsdir ./srsdir run --group all --chainId 11155111 --zkbtcBridgeAddr 0x49793ff075b696e6bef1b85e4e85fe669041b312 --icpPublickey 023f203422be55a3576f46dc6770bdc7865a126381c1963a2d82b49f4158409a2e

    ./setup  --datadir ./testdir --srsdir ./srsdir run --group ethereum --chainId 11155111 --zkbtcBridgeAddr 0x49793ff075b696e6bef1b85e4e85fe669041b312 

    ./setup  --datadir ./testdir --srsdir ./srsdir run --group bitcoin --icpPublickey 023f203422be55a3576f46dc6770bdc7865a126381c1963a2d82b49f4158409a2e




    ./setup readCircuit readCircuit --file /opt/lightec/opt2/circuitsetup/beta3/redeem.vk


    ./setup  --datadir ./testdir --srsdir ./srsdir run --type ethTxInEth2 --chainId 11155111 --zkbtcBridgeAddr 0xA986b6Ae23Da8c9d06074033BDb6B3e0421De346 
