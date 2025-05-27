

## Import bitcoin chain proof

blockchain

    ./node proof import --proof /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/blockchain_0_901151.proof --witness /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/blockchain_0_901151.wtns --name btcDuperRecursive_899136_901152  --datadir /Users/red/lworkspace/lightec/audit/daemon/testdata/.daemon


## Import ethereum chain proof

outer

     ./node proof import --proof /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1471_outer.proof --witness /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1471_outer.wtns --name outer_1471  --datadir /Users/red/lworkspace/lightec/audit/daemon/testdata/.daemon

unit

    ./node proof import --proof /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1471_unit.proof --witness /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1471_unit.wtns --name unit_1471  --datadir /Users/red/lworkspace/lightec/audit/daemon/testdata/.daemon

duty

    ./node proof import --proof /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1471_duty.proof --witness /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1471_duty.wtns --name duty_1471  --datadir /Users/red/lworkspace/lightec/audit/daemon/testdata/.daemon

recursive

    ./node proof import --proof /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1472_recursive.proof --witness /Users/red/lworkspace/lightec/audit/daemon/testdata/mainnet/sc1472_recursive.wtns --name recursive_1471  --datadir /Users/red/lworkspace/lightec/audit/daemon/testdata/.daemon
