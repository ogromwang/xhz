package model

type Result struct {
	Code  uint32      `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

type ResultWithPage struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}

type ResultWithMore struct {
	List interface{} `json:"list"`
	More bool        `json:"more"`
}
