package models

type Nmap_Payload struct {
	Hosts []Host `json:"hosts"`
	Verbose string `json:"verbose"`
}