package models

type John_Result struct {
	Passwords map[string]string `json:"crckdPswds"`
	AvailableFiles []string `json:"availableFiles"`
	Details string `json:"details"`
	Error string `json:"error"`
	Result string `json:"result"`
}