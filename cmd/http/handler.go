package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"errors"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"

	"github.com/craguilar/event-management-service/internal/app"
	"github.com/gorilla/mux"
)

type EventServiceHandler struct {
	eventService app.EventService
	guestService app.GuestService
	taskService  app.TaskService
}

func NewServiceHandler(event app.EventService, guest app.GuestService, task app.TaskService) *EventServiceHandler {
	return &EventServiceHandler{
		eventService: event,
		guestService: guest,
		taskService:  task,
	}
}

func (c *EventServiceHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	// Check associated user
	user, err := getUser(r)
	if err != nil {
		log.Warn("Error when decoding Authorization ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Invalid Authorization header"))
		return
	}
	// Body decode
	var event app.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Warn("Error when decoding Body", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Invalid Body parameter"))
		return
	}
	// And finally create the user
	createdCar, err := c.eventService.CreateOrUpdate(user, &event)
	if err != nil {
		log.Error("Error when creating event ", err)
		WriteError(w, http.StatusInternalServerError, err)
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
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if car == nil {
		WriteError(w, http.StatusNotFound, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(car))
}

func (c *EventServiceHandler) ListEvent(w http.ResponseWriter, r *http.Request) {
	// Check associated user
	user, err := getUser(r)
	if err != nil {
		log.Warn("Error when decoding Authorization ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Invalid Authorization header"))
		return
	}
	// Then list
	events, err := c.eventService.List(user)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(events))
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
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Tasks

func (c *EventServiceHandler) AddTask(w http.ResponseWriter, r *http.Request) {

	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expected eventId")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expectde eventId as query parameter"))
		return
	}

	var task app.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Warn("Error when decoding Body", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Invalid Body parameter"))
		return
	}
	createdTask, err := c.taskService.CreateOrUpdate(eventId, &task)
	if err != nil {
		log.Error("Error when creating task ", err)
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(createdTask))
}

func (c *EventServiceHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expected eventId")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expected eventId as query parameter"))
		return
	}
	vars := mux.Vars(r)
	taskId, ok := vars["taskId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "BadRequest"))
		return
	}
	task, err := c.taskService.Get(eventId, taskId)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if task == nil {
		WriteError(w, http.StatusNotFound, nil)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(task))
}

func (c *EventServiceHandler) ListTask(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warnf("Expected eventId got %s", eventId)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expected eventId as query parameter"))
		return
	}
	tasks, err := c.taskService.List(eventId)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(tasks))
}

func (c *EventServiceHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expected eventId")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "Expected eventId as query parameter"))
		return
	}

	vars := mux.Vars(r)
	taskId, ok := vars["taskId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(SerializeError(http.StatusBadRequest, "BadRequest"))
		return
	}
	err := c.taskService.Delete(eventId, taskId)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *EventServiceHandler) AddGuest(w http.ResponseWriter, r *http.Request) {
	var guest app.Guest
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expected eventId")
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
		log.Error("Error when creating event ", err)
		WriteError(w, http.StatusInternalServerError, err)
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
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	if car == nil {
		WriteError(w, http.StatusNotFound, nil)
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
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(SerializeData(guests))
}

func (c *EventServiceHandler) DeleteGuest(w http.ResponseWriter, r *http.Request) {
	eventId := r.URL.Query().Get("eventId")
	if eventId == "" {
		log.Warn("Expected eventId")
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
		WriteError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getUser(r *http.Request) (string, error) {
	authorization := r.Header.Get("Authorization")
	authDetails := strings.Split(authorization, " ")
	if len(authDetails) < 2 {
		return "", errors.New("invalid authorization info")
	}
	// Handle the use case of local testing
	if IsLocal() && strings.HasPrefix(authDetails[1], "dummy") {
		return "dummy", nil
	}
	parser := new(jwt.Parser)
	tokenString := authDetails[1]
	claims := jwt.MapClaims{}
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return "", err
	}
	return claims["email"].(string), nil
}
