package app

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// TODO on 01-20-2023: Add pagination and stuff too busy for it now

type AuthorizationService interface {
	Authorize(userName, eventId string) bool
}

// Interface for Event service
type EventService interface {
	Get(id string) (*Event, error)
	List(eventManager string) ([]*EventSummary, error)
	ListBy(filter func(*EventSummary) bool) ([]*EventSummary, error)
	ListOwners(id string) (*EventSharedEmails, error)
	CreateOrUpdate(eventManager string, u *Event) (*Event, error)
	CreateOwner(eventManager string, u *EventSharedEmails) (*EventSharedEmails, error)
	Delete(eventManager, id string) error
}

type EventActions interface {
	// Filter based on the referenced function
	SendPendingTasksNotifications() error
}

type GuestService interface {
	Get(eventManager, eventId, id string) (*Guest, error)
	List(eventManager, eventId string) ([]*Guest, error)
	CopyFrom(eventManager string, eventId string, copy *CopyGuestRequest) error
	CreateOrUpdate(eventManager, eventId string, u *Guest) (*Guest, error)
	Delete(eventManager, eventId, id string) error
}

type TaskService interface {
	Get(eventId, id string) (*Task, error)
	List(eventId string) ([]*Task, error)
	CreateOrUpdate(eventId string, u *Task) (*Task, error)
	Delete(eventId, id string) error
}

type ExpenseService interface {
	Get(eventId, id string) (*ExpenseCategory, error)
	List(eventId string) ([]*ExpenseCategory, error)
	CreateOrUpdate(eventId string, u *ExpenseCategory) (*ExpenseCategory, error)
	Delete(eventId, id string) error
}

// Event : Required Name , MainLocation, EventDay. An event has Guests ,Expenses and Tasks
type Event struct {
	Id                  string    `json:"id"`
	Name                string    `json:"name" validate:"required"`
	MainLocation        string    `json:"mainLocation" validate:"required"`
	EventDay            time.Time `json:"eventDay" validate:"required"`
	Description         string    `json:"description"`
	Guests              []*Guest  `json:"guests"`
	NotificationEnabled bool      `json:"isNotificationEnabled"`
	v                   *validator.Validate
	TimeCreatedOn       time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn       time.Time `json:"timeUpdatedOn"`
}

type EventSummary struct {
	Id                  string    `json:"id" validate:"required"`
	Name                string    `json:"name" validate:"required"`
	MainLocation        string    `json:"mainLocation" validate:"required"`
	EventDay            time.Time `json:"eventDay" validate:"required"`
	NotificationEnabled bool      `json:"isNotificationEnabled"`
	TimeCreatedOn       time.Time `json:"timeCreatedOn"`
}

type EventOwner struct {
	OwnerEmail   string `json:"ownerEmail" validate:"required"`
	EventSummary *EventSummary
}

type EventSharedEmails struct {
	EventId      string   `json:"eventId" validate:"required"`
	SharedEmails []string `json:"sharedEmails" validate:"required"`
}

// Guest : Required FirstName,LastName,Tentative,NumberOfSeats
type Guest struct {
	Id             string `json:"id"`
	FirstName      string `json:"firstName" validate:"required"`
	LastName       string `json:"lastName" validate:"required"`
	GuestOf        string `json:"guestOf"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Tentative      bool   `json:"isTentative"`
	Country        string `json:"country"`
	State          string `json:"state"`
	RequiresInvite bool   `json:"requiresInvite"`
	NotAttending   bool   `json:"isNotAttending"`
	NumberOfSeats  int    `json:"numberOfSeats" validate:"required"`
	v              *validator.Validate
	TimeCreatedOn  time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn  time.Time `json:"timeUpdatedOn"`
}

type CopyGuestRequest struct {
	FromEvent string `json:"fromEvent"`
}

type Task struct {
	Id            string `json:"id"`
	Name          string `json:"name" validate:"required"`
	Status        string `json:"status" validate:"required"` // PENDING, DONE
	v             *validator.Validate
	TimeCreatedOn time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn time.Time `json:"timeUpdatedOn"`
}

// Expense representation , the Category MUST be unique per eventId
type ExpenseCategory struct {
	Id              string     `json:"id"`
	Category        string     `json:"category" validate:"required"`
	AmountProjected float64    `json:"amountProjected"`
	AmountPaid      float64    `json:"amountPaid"`
	AmountTotal     float64    `json:"amountTotal"`
	Expenses        []*Expense `json:"expenses"`
	v               *validator.Validate
	TimeCreatedOn   time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn   time.Time `json:"timeUpdatedOn"`
}

type Expense struct {
	Id            string    `json:"id"`
	WhoPaid       string    `json:"whoPaid" validate:"required"`
	TimePaidOn    time.Time `json:"timePaidOn"`
	AmountPaid    float64   `json:"amountPaid"`
	v             *validator.Validate
	TimeCreatedOn time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn time.Time `json:"timeUpdatedOn"`
}

func (e *Event) Validate() error {
	if e.v == nil {
		e.v = validator.New()
	}
	return e.v.Struct(e)
}

func (e *Event) ToSummary() *EventSummary {
	return &EventSummary{
		Id:                  e.Id,
		Name:                e.Name,
		MainLocation:        e.MainLocation,
		EventDay:            e.EventDay,
		TimeCreatedOn:       e.TimeCreatedOn,
		NotificationEnabled: e.NotificationEnabled,
	}
}

func (e *ExpenseCategory) Validate() error {
	if e.v == nil {
		e.v = validator.New()
	}
	// Arbitrart number to avoid ExpenseCategory row growing large
	if len(e.Expenses) > 80 {
		return errors.New("trying to create more than 80 Expenses  . Not allowed")
	}

	for _, value := range e.Expenses {
		if err := value.Validate(); err != nil {
			return err
		}
	}
	return e.v.Struct(e)
}

func (e *Expense) Validate() error {
	if e.v == nil {
		e.v = validator.New()
	}
	return e.v.Struct(e)
}

func (t *Task) Validate() error {
	if t.v == nil {
		t.v = validator.New()
	}
	return t.v.Struct(t)
}

func (g *Guest) Validate() error {
	if g.v == nil {
		g.v = validator.New()
	}
	return g.v.Struct(g)
}

func GenerateId(value string) string {
	data := []byte(strings.ToUpper(value))
	return fmt.Sprintf("%x", md5.Sum(data))
}

func GenerateRandomId() (string, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
