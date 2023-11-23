package bitcoin

const (
	// blockchain
	GETBLOCKHEADER = "getblockheader"
	GETBLOCKHASH   = "getblockhash"
	GETBLOCKCOUNT  = "getblockcount"
	SCANTXOUTSET   = "scantxoutset"

	// utils
	CREATEMULTISIG = "createmultisig"

	//wallet
	GETRAWCHANGEADDRESS = "getrawchangeaddress"
	DUMPPRIVKEY         = "dumpprivkey"
	GETADDRESSINFO      = "getaddressinfo"

	//rawtransaction
	CREATERAWTRANSACTION      = "createrawtransaction"
	SIGNRAWTRANSACTIONWITHKEY = "signrawtransactionwithkey"
	SENDRAWTRANSACTION        = "sendrawtransaction"
	GETRAWTRANSACTION         = "getrawtransaction"
)

type AddrType string

const (
	BECH32M AddrType = "bech32m"
	BECH32  AddrType = "bech32"
)
