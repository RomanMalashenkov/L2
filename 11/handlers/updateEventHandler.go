package handlers

import (
	"L2/11/cache"
	"L2/11/domain"
	"L2/11/models"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func UpdateEventHandler(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	if r.Method != http.MethodPost {
		domain.ErrorLogger(w, errors.New("method error"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var decoded models.UpdEvent

	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		domain.ErrorLogger(w, decodingBodyErr)
		return
	}

	dateQuery := decoded.Date
	timeQuery := decoded.Time
	userQuery := decoded.UserId
	newUserQuery := decoded.NewData

	if _, errParse := time.Parse("2006-01-02", dateQuery); errParse != nil {
		domain.ErrorLogger(w, errParse)
		return
	}

	if _, errParseNewDate := time.Parse("2006-01-02", newUserQuery.Date); errParseNewDate != nil {
		domain.ErrorLogger(w, errParseNewDate)
		return
	}

	if _, errParseTime := time.Parse("15:00", timeQuery); errParseTime != nil {
		domain.ErrorLogger(w, errParseTime)
		return
	}

	if _, errParseNewTime := time.Parse("15:00", newUserQuery.Time); errParseNewTime != nil {
		domain.ErrorLogger(w, errParseNewTime)
		return
	}

	c.Update(models.Event{Date: dateQuery, Time: timeQuery, UserId: userQuery}, newUserQuery.Date, newUserQuery.Time)

	domain.ResponseLogger(w, "event updated")
}
