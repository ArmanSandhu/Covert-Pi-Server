package models

type Airmon_Result struct {
	Airmon_Interfaces []Airmon_Interface `json:"airmon_interfaces"`
	Details string `json:"details"`
	Result string `json:"result"`
}