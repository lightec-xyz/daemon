# Generator

## Build

    makdir <projectDir>

    cd <projectDir>

    git clone https://<Token>@github.com/lightec-xyz/provers.git
    cd provers
    git checkout v0.2.1

    git cloen https://<Token>@github.com/lightec-xyz/btc_provers.git
    cd btc_provers
    git checkout v0.2.0  

    git clone https://<Token>@github.com/lightec-xyz/reLight.git
    cd reLight
    git checkout v0.2.0  


    git clone https://<Token>@github.com/lightec-xyz/daemon
    cd daemon
    git checkout v0.2.1
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

update go version to v1.22.3 


## Run

    ./generator --config ./testnet.json run


## Setup circuits locally
如果很难 从mac studio下载ccs/pk/vk，可以采用local setup的方式生成ccs/pk/vk/sol 文件。以下是操作步骤。

1. 从mac studio 下载aztec_bn254的srs/lsrs到本地
```
 scp -r <userName>@58.41.9.129:/opt/lightec/srs/aztec_bn254  <本地保存srs的路径>
```

2. setup provers中的各个模块

2.1 下载并编译 provers中的各个模块

```
git clone https://<Token>@github.com/lightec-xyz/provers.git
cd provers 
git checkout v0.2.0
```

copy下面的脚本 放在provers 目录下并执行， 编译各个模块
```
#!/bin/bash

go mod tidy

cd cmd/beacon-header
go build

cd ../beacon-header-finality
go build

cd ../tx-in-eth2
go build

cd ../redeem
go build 

```

2.2 setup prover 中的各个模块
将cmd/beacon-header/setup.sh, cmd/beacon-header-finality/setup.sh, cmd/tx-in-eth2/setup.sh, cmd/redeem/setup.sh 中的home 和 srs 修改为本机对应的目录。

例如：
```
export home=/Users/xxx/lightec-xyz/worker/circuit_data/circuits  # tx-in-eth2/setup.sh 中的是 "data"
export srs=/Users/xxx/lightec-xyz/srs/aztec_bn254
```

依次进入md/beacon-header，cmd/beacon-header-finality， cmd/beacon-header-finality， cmd/redeem 目录并运行
```
sh setup.sh # 注意tx-in-eth2的setup 采用的nohup的运行方式，需要Ctrl+C 退出。
```

3. setup reLight中的各个模块

3.1 下载并编译 reLight中的各个模块

```
  git clone https://<Token>@github.com/lightec-xyz/reLight.git
  cd reLight
  git checkout v0.2.0  
```

copy下面的脚本 放在reLight目录下并执行， 编译各个模块
```
#!/bin/bash

go mod tidy

cd cmd/unit
go build

cd ../genesis
go build

cd ../recursive
go build

```
3.2 setup reLight 中的各个模块
将cmd/unit/setup.sh, cmd/genesis/setup.sh, cmd/recursive/setup.sh  中的data 和 srs 修改为本机对应的目录。

例如：
```
export data=/Users/xxx/lightec-xyz/worker/circuit_data/circuits
export srs=/Users/xxx/lightec-xyz/srs/aztec_bn254
```

将cmd/unit/setup.sh 中的 "./unit sc"这一行注释掉。 

```
#./unit sc --datadir $data --subdir "../data/sc/unit/sc181" --file holesky_sync_committee_update_181.json
```

依次进入cmd/unit cmd/genesis cmd/recursive 目录并运行
```
sh setup.sh 
```

3. setup btc_prover中的各个模块(TBD）， 因为btc 的ccs/pk/vk文件只有grandrollup.ccs/pk/vk，建议直接从mac studio:/opt/lightec/circuit_data/beta1/circuits/grandrollup 中下载。
   
   




