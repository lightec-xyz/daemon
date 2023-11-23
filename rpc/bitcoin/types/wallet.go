package types

type AddressInfo struct {
	Address             string   `json:"address"`
	Desc                string   `json:"desc"`
	Hdkeypath           string   `json:"hdkeypath"`
	Hdmasterfingerprint string   `json:"hdmasterfingerprint"`
	Hdseedid            string   `json:"hdseedid"`
	Ischange            bool     `json:"ischange"`
	Ismine              bool     `json:"ismine"`
	Isscript            bool     `json:"isscript"`
	Iswatchonly         bool     `json:"iswatchonly"`
	Iswitness           bool     `json:"iswitness"`
	Labels              []string `json:"labels"`
	Pubkey              string   `json:"pubkey"`
	ScriptPubKey        string   `json:"scriptPubKey"`
	Solvable            bool     `json:"solvable"`
	Timestamp           int      `json:"timestamp"`
	WitnessProgram      string   `json:"witness_program"`
	WitnessVersion      int      `json:"witness_version"`
}
