package main

import (
	"L2/11/cache"
	"L2/11/handlers"
	"log"
	"net/http"
)

func main() {

	storage := cache.NewCache()

	mux := http.NewServeMux()

	createEventHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.CreateEventHandler(writer, request, storage)
	})

	updateEventHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.UpdateEventHandler(writer, request, storage)
	})

	deleteEventHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.DeleteEventHandler(writer, request, storage)
	})

	getDayEventHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.GetDayEventHandler(writer, request, storage)
	})

	getWeekEventHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.GetWeekEventHandler(writer, request, storage)
	})

	getMonthEventHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handlers.GetMonthEventHandler(writer, request, storage)
	})

	mux.Handle("/create_event", mwLogger(createEventHandler))
	mux.Handle("/update_event", mwLogger(updateEventHandler))
	mux.Handle("/delete_event", mwLogger(deleteEventHandler))
	mux.Handle("/events_for_day", mwLogger(getDayEventHandler))
	mux.Handle("/events_for_week", mwLogger(getWeekEventHandler))
	mux.Handle("/events_for_month", mwLogger(getMonthEventHandler))

	log.Println("Server started on localhost:8000...")

	log.Fatal(http.ListenAndServe(":8000", mux))
}

func mwLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Executing %s...", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
