package app

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// TODO on 01-20-2023: Add pagination and stuff too busy for it now

// Interface for Event service
type EventService interface {
	Get(id string) (*Event, error)
	List(user string) ([]*EventSummary, error)
	CreateOrUpdate(u *Event) (*Event, error)
	Delete(id string) error
}

type GuestService interface {
	Get(eventId, id string) (*Guest, error)
	List(eventId string) ([]*Guest, error)
	CreateOrUpdate(eventId string, u *Guest) (*Guest, error)
	Delete(eventId, id string) error
}

type LocationService interface {
	Get(eventId, id string) (*Location, error)
	List(eventId string) ([]*Location, error)
	CreateOrUpdate(u *Location) (*Location, error)
	Delete(eventId, id string) error
}

type TaskService interface {
	Get(eventId, id string) (*Guest, error)
	List(eventId string) ([]*Guest, error)
	CreateOrUpdate(u *Guest) (*Guest, error)
	Delete(eventId, id string) error
}

type ExpenseService interface {
	Get(eventId, id string) (*Expense, error)
	List(eventId string) ([]*Expense, error)
	CreateOrUpdate(u *Expense) (*Expense, error)
	Delete(eventId, id string) error
}

type EventSummary struct {
	Id            string    `json:"id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	MainLocation  string    `json:"mainLocation" validate:"required"`
	EventDay      time.Time `json:"eventDay" validate:"required"`
	TimeCreatedOn time.Time `json:"timeCreatedOn"`
}

// Event : Required Name , MainLocation, EventDay
type Event struct {
	Id            string      `json:"id"`
	Name          string      `json:"name" validate:"required"`
	MainLocation  string      `json:"mainLocation" validate:"required"`
	EventDay      time.Time   `json:"eventDay" validate:"required"`
	Description   string      `json:"description"`
	Expenses      []*Expense  `json:"expenses"`
	Guests        []*Guest    `json:"guests"`
	Locations     []*Location `json:"locations"`
	Tasks         []*Task     `json:"tasks"`
	v             *validator.Validate
	TimeCreatedOn time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn time.Time `json:"timeUpdatedOn"`
}

// Guest : Required FirstName,LastName,Tentative,NumberOfSeats
type Guest struct {
	Id            string `json:"id"`
	FirstName     string `json:"firstName" validate:"required"`
	LastName      string `json:"lastName" validate:"required"`
	GuestOf       string `json:"guestOf"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Tentative     bool   `json:"isTentative"`
	Country       string `json:"country"`
	State         string `json:"state"`
	NumberOfSeats int    `json:"numberOfSeats" validate:"required"`
	v             *validator.Validate
	TimeCreatedOn time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn time.Time `json:"timeUpdatedOn"`
}

type Location struct {
	Id            uuid.UUID `json:"id" validate:"required"`
	Name          string    `json:"name"`
	Where         string    `json:"where" validate:"required"`
	When          time.Time `json:"when" validate:"required"`
	v             *validator.Validate
	TimeCreatedOn time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn time.Time `json:"timeUpdatedOn"`
}

type Expense struct {
	Id              uuid.UUID `json:"id" validate:"required"`
	Name            string    `json:"name" validate:"required"`
	AmountProjected float64   `json:"amountProjected"`
	AmountActual    float64   `json:"amountActual"`
	AmountPaid      float64   `json:"amountPaid"`
	v               *validator.Validate
	TimeCreatedOn   time.Time `json:"timeCreatedOn"`
	TimeUpdatedOn   time.Time `json:"timeUpdatedOn"`
}

type Task struct {
	Id            uuid.UUID `json:"id" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Status        string    `json:"status" validate:"required"` // PENDING, DONE
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

func (l *Location) Validate() error {
	if l.v == nil {
		l.v = validator.New()
	}
	return l.v.Struct(l)
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