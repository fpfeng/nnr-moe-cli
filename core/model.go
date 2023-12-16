package core

type ServerInfo struct {
	Sid    string   `json:"sid"`
	Name   string   `json:"name"`
	Host   string   `json:"host"`
	Min    int      `json:"min"`
	Max    int      `json:"max"`
	Mf     float64  `json:"mf"`
	Level  int      `json:"level"`
	Detail string   `json:"detail"`
	Types  []string `json:"types"`
}

type ServerList []ServerInfo

type ResponseBase struct {
	Status int32
}

func (rb *ResponseBase) IsStatusOK() bool {
	return rb.Status == 1
}

type ResponseServers struct {
	ResponseBase
	Data ServerList
}

type RuleBase struct {
	Remote  string `json:"remote"`
	Rport   int    `json:"rport"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Setting struct {
		LoadbalanceMode string `json:"loadbalanceMode"`
	} `json:"setting"`
}

type Rule struct {
	RuleBase
	Rid     string `json:"rid"`
	Uid     string `json:"uid"`
	Sid     string `json:"sid"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Status  int    `json:"status"`
	Traffic int    `json:"traffic"`
	Data    string `json:"data"`
	Date    int64  `json:"date"`
	Setting struct {
		ProxyProtocol   int           `json:"proxyProtocol"`
		LoadbalanceMode string        `json:"loadbalanceMode"`
		Mix0Rtt         bool          `json:"mix0rtt"`
		Src             interface{}   `json:"src"`
		Cfips           []interface{} `json:"cfips"`
	} `json:"setting"`
}

type RuleList []Rule

type ResponseRules struct {
	ResponseBase
	Data RuleList
}

type RequestAddRule struct {
	Sid string `json:"sid"`
	RuleBase
}

type ResponseRuleDetail struct {
	ResponseBase
	Data Rule
}

type RequestRuleRid struct {
	Rid string `json:"rid"`
}

type ResponseDeleteRule struct {
	ResponseBase
	Data RequestRuleRid
}

type RequestEditedRule struct {
	Rid string `json:"rid"`
	RuleBase
}

type TypeRequestBody interface {
	*RequestAddRule | *RequestRuleRid | *RequestEditedRule | string
}

type TypeResponse interface {
	*ResponseServers | *ResponseRuleDetail | *ResponseDeleteRule | *ResponseRules
	IsStatusOK() bool
}
