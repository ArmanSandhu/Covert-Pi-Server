package models

type Airodump_AP struct {
	Mac string `json:"apmac"`
	First_Seen string `json:"apfirsttimeseen"`
	Last_Seen string `json:"aplasttimeseen"`
	Channel string `json:"apchannel"`
	Speed string `json:"apspeed"`
	Privacy string `json:"apprivacy"`
	Cipher string `json:"apcipher"`
	Auth string `json:"apauth"`
	Power string `json:"appower"`
	Num_Packets string `json:"apnumpackets"`
	IV string `json:"apiv"`
	Lan_IP string `json:"aplanip"`
	ID_Len string `json:"apidlen"`
	Essid string `json:"apessid"`
	Key string `json:"apkey"`
	Clients []Airodump_Client `json:"apclients"`
}