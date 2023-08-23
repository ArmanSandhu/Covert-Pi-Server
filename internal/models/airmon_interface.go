package models

type Airmon_Interface struct {
	PHY string `json:"phy"`
	Interface string `json:"interface"`
	Driver string `json:"driver"`
	Chipset string `json:"chipset"`
	Mode string `json:"mode"`
}