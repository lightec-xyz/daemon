# API doc

API documentation for interfacing with the zkBTC daemon, similar to Ethereum integration workflows.

### common api

zkbtc_version

     curl -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"'"zkbtc_version"'","params":[],"id":1}' http://127.0.0.1:9977

you can find more api interfaces [here](https://github.com/lightec-xyz/daemon/blob/74148d4cc671786a909defad26b3216fc4aa7102/rpc/api.go#L7)

### admin api

some special apis require an admin role [more admin api](https://github.com/lightec-xyz/daemon/blob/74148d4cc671786a909defad26b3216fc4aa7102/rpc/api.go#L63).

    curl -X POST -H "Content-Type: application/json" -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdWJsaWNLZXkiOjIsImV4cCI6MTc0OTU2NTY3NX0.MBBAptLkmG9spEhznOUPRGbJmctnINkbZHuhVHNzED4" -d '{"jsonrpc":"2.0","method":"zkbtc_removeUnGenProof","params":["dfccd240d664b642c6017064566c5e28b9b6309b27f9794ae48dfd59a6ed8216"],"id":1}' http://127.0.0.1:9977

you can generate the jwt token using [source code](https://github.com/lightec-xyz/daemon/blob/74148d4cc671786a909defad26b3216fc4aa7102/rpc/jwt_test.go#L11) or the cmd  `./node --config ./local.json jwt`


autoSubmitMaxValue

    curl -X POST -H "Content-Type: application/json"  -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdWJsaWNLZXkiOjIsImV4cCI6MTc0OTU2NTY3NX0.MBBAptLkmG9spEhznOUPRGbJmctnINkbZHuhVHNzED4"  -d '{"jsonrpc":"2.0","method":"'"zkbtc_autoSubmitMaxValue"'","params":[100000000],"id":1}' http://127.0.0.1:9977

autoSubmitMinValue

    curl -X POST -H "Content-Type: application/json"  -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdWJsaWNLZXkiOjIsImV4cCI6MTc0OTU2NTY3NX0.MBBAptLkmG9spEhznOUPRGbJmctnINkbZHuhVHNzED4"  -d '{"jsonrpc":"2.0","method":"'"zkbtc_autoSubmitMinValue"'","params":[21000],"id":1}' http://127.0.0.1:9977

setGasPriceLimit

    curl -X POST -H "Content-Type: application/json"  -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdWJsaWNLZXkiOjIsImV4cCI6MTc0OTU2NTY3NX0.MBBAptLkmG9spEhznOUPRGbJmctnINkbZHuhVHNzED4"  -d '{"jsonrpc":"2.0","method":"'"zkbtc_setGasPrice"'","params":[3000000000],"id":1}' http://127.0.0.1:9977
