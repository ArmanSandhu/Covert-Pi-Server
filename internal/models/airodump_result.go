package models

type Airodump_Result struct {
	Airodump_APs []Airodump_AP `json:"airodump_aps"`
	Message string `json:"message"`
	Result string `json:"result"`
}