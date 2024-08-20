package main

type RelayInfo struct {
	Contact        string `json:"contact"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	Pubkey         string `json:"pubkey"`
	Software       string `json:"software"`
	Supported_nips []int  `json:"supported_nips"`
	Version        string `json:"version"`
}
