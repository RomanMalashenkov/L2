package models

type Event struct {
	UserId int    `json:"user_id"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Uid    int64  `json:"uid"`
}
