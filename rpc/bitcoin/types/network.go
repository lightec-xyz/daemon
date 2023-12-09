package types

type NetworkInfo struct {
	Connections        int           `json:"connections"`
	ConnectionsIn      int           `json:"connections_in"`
	ConnectionsOut     int           `json:"connections_out"`
	Incrementalfee     float64       `json:"incrementalfee"`
	Localaddresses     []interface{} `json:"localaddresses"`
	Localrelay         bool          `json:"localrelay"`
	Localservices      string        `json:"localservices"`
	Localservicesnames []string      `json:"localservicesnames"`
	Networkactive      bool          `json:"networkactive"`
	Networks           []struct {
		Limited                   bool   `json:"limited"`
		Name                      string `json:"name"`
		Proxy                     string `json:"proxy"`
		ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
		Reachable                 bool   `json:"reachable"`
	} `json:"networks"`
	Protocolversion int     `json:"protocolversion"`
	Relayfee        float64 `json:"relayfee"`
	Subversion      string  `json:"subversion"`
	Timeoffset      int     `json:"timeoffset"`
	Version         int     `json:"version"`
	Warnings        string  `json:"warnings"`
}
