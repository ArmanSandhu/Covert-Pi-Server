package models

type Port struct {
	Num string `json:"portnum"`
	Porttype string `json:"porttype"`
	State string `json:"portstate"`
	Service string `json:"portservice"`
}