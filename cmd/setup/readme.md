# circuit Setup

#### Command

generate single circuit setup:(beaconInner, beaconOuter, beaconUnit, beaconGenesis, beaconRecursive,ethTxInEth2,ethBeaconHeader,
ethFinalityHeader,ethRedeem, btcBase, btcMiddle, btcUpper,btcBlockChain,btcBlockDepth,btcTxInChain)

    ./setup  --datadir ./testdir --srsdir ./srsdir run --type beaconOuter 

batch generate circuit setup:(all, bitcoin, beacon, ethereum)

    ./setup  --datadir ./testdir --srsdir ./srsdir run --group all 