package models

type Airodump_Client struct {
	Mac string `json:"clientmac"`
	First_Seen string `json:"clientfirsttimeseen"`
	Last_Seen string `json:"clientlasttimeseen"`
	Power string `json:"clientpower"`
	Num_Packets string `json:"clientnumpackets"`
	Bssid string `json:"clientbssid"`
	Probed_Essid string `json:"clientprobedessid"`
}