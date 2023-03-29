package models

type Host struct {
	Name string `json:"hostname"`
	Ip string `json:"hostip"`
	Status string `json:"hoststatus"`
	Latency string `json:"hostlatency"`
	Ports []Port `json:"hostports"`
	Mac string `json:"hostmac"`
}