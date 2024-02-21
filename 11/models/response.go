package models

type Response struct {
	Body Result `json:"result"`
}

type Result struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}
