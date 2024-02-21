package models

type Error struct {
	Err Details `json:"error"`
}

type Details struct {
	ErrCode    int    `json:"code"`
	ErrMessage string `json:"message"`
}
