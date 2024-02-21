package domain

import (
	"L2/11/models"
	"encoding/json"
	"log"
	"net/http"
	time2 "time"
)

func NewEvent(date, time string, userId int) *models.Event {
	Uid := time2.Now().Unix()
	return &models.Event{
		UserId: userId,
		Date:   date,
		Time:   time,
		Uid:    Uid,
	}
}

func ErrorLogger(w http.ResponseWriter, errorInput error) {
	log.Println(errorInput)
	w.WriteHeader(http.StatusInternalServerError)
	details := models.Details{ErrCode: http.StatusInternalServerError, ErrMessage: errorInput.Error()}
	response, _ := json.MarshalIndent(models.Error{Err: details}, "", "\t")
	w.Write(response)
	return
}

func ResponseLogger(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusInternalServerError)
	result := models.Result{StatusCode: http.StatusOK, Message: message}
	response, _ := json.MarshalIndent(models.Response{Body: result}, "", "\t")
	w.Write(response)
	return
}
