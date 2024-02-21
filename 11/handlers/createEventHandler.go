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

func CreateEventHandler(w http.ResponseWriter, r *http.Request, c *cache.Cache) {
	if r.Method != http.MethodPost {
		domain.ErrorLogger(w, errors.New("method error"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var decoded models.Event

	if decodingBodyErr := decoder.Decode(&decoded); decodingBodyErr != nil {
		domain.ErrorLogger(w, decodingBodyErr)
		return
	}

	dateQuery := decoded.Date
	timeQuery := decoded.Time
	userQuery := decoded.UserId

	if _, errParse := time.Parse("2006-01-02", dateQuery); errParse != nil {
		domain.ErrorLogger(w, errParse)
		return
	}

	if _, errParseTime := time.Parse("15:00", timeQuery); errParseTime != nil {
		domain.ErrorLogger(w, errParseTime)
		return
	}

	event := domain.NewEvent(dateQuery, timeQuery, userQuery)
	c.Create(event)

	domain.ResponseLogger(w, "event created")
}
