package http

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/craguilar/event-management-service/internal/app"
	"github.com/craguilar/event-management-service/internal/app/mock"
	"github.com/gorilla/mux"
)

type EventServiceHandler struct {
	eventService app.EventService
	guestService app.GuestService
}

func NewServiceHandler(service app.EventService) *EventServiceHandler {
	return &EventServiceHandler{
		eventService: service,
		guestService: mock.NewGuestService(service),
	}
}

func (c *EventServiceHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	var event app.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Warn("Error when decoding Body", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Invalid Body parameter"))
		return
	}
	createdCar, err := c.eventService.CreateOrUpdate(&event)
	if err != nil {
		log.Error("Error when creating car", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(createdCar))
}

func (c *EventServiceHandler) GetEvent(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	plate, ok := vars["eventId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "BadRequest"))
		return
	}
	car, err := c.eventService.Get(plate)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if car == nil {
		writeError(w, http.StatusNotFound, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(car))
}

func (c *EventServiceHandler) ListEvent(w http.ResponseWriter, r *http.Request) {
	cars, err := c.eventService.List()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(cars))
}

func (c *EventServiceHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventId, ok := vars["eventId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "BadRequest"))
		return
	}
	err := c.eventService.Delete(eventId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Guests

func (c *EventServiceHandler) AddGuest(w http.ResponseWriter, r *http.Request) {
	var guest app.Guest
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expectde eventId")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expectde eventId as query parameter"))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&guest)
	if err != nil {
		log.Warn("Error when decoding Body", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Invalid Body parameter"))
		return
	}
	createdGuest, err := c.guestService.CreateOrUpdate(eventId, &guest)
	if err != nil {
		log.Error("Error when creating car", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(createdGuest))
}

func (c *EventServiceHandler) GetGuest(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expectde eventId")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expectde eventId as query parameter"))
		return
	}
	vars := mux.Vars(r)
	guestId, ok := vars["guestId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "BadRequest"))
		return
	}
	car, err := c.guestService.Get(eventId, guestId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if car == nil {
		writeError(w, http.StatusNotFound, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(car))
}

func (c *EventServiceHandler) ListGuest(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warnf("Expected eventId got %s", eventId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expected eventId as query parameter"))
		return
	}
	guests, err := c.guestService.List(eventId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(guests))
}

func (c *EventServiceHandler) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expectde eventId")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expectde eventId as query parameter"))
		return
	}

	vars := mux.Vars(r)
	guestId, ok := vars["guestId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "BadRequest"))
		return
	}
	err := c.guestService.Delete(eventId, guestId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Write HTTP Error out to http.ResponseWriter
func writeError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	var errorCode string
	if statusCode == 400 {
		errorCode = "InvalidParameter."
	} else if statusCode == 404 {
		errorCode = "NotFound or caller don't have access."
	} else if statusCode == 401 {
		errorCode = "Unauthorized"
	} else if statusCode == 409 {
		errorCode = "Conflict with resource"
	} else if statusCode == 500 {
		errorCode = "InternalServerError"
	}
	log.Warnf("Received error from call %s", err)
	w.Write(SerializeError(statusCode, errorCode))

}
