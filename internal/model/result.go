package model

type Result struct {
	Code uint32 `json:"code"`
	Data interface{} `json:"data"`
	Error string `json:"error"`
}