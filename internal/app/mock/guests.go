package mock

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/craguilar/event-management-service/internal/app"
	"golang.org/x/exp/slices"
)

type GuestService struct {
	eventService app.EventService
	lock         sync.RWMutex
}

func NewGuestService(eventService app.EventService) *GuestService {
	return &GuestService{
		eventService: eventService,
	}
}

func (c *GuestService) Get(eventId, id string) (*app.Guest, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	// Get the value
	value, err := c.eventService.Get(eventId)
	if err != nil {
		return nil, nil
	}
	guest, _, err := searchByGuestId(value.Guests, id)
	return guest, err
}

func (c *GuestService) CopyFrom(userName string, eventId string, copy *app.CopyGuestRequest) error {
	return errors.New("not implemented")
}

func (c *GuestService) List(eventId string) ([]*app.Guest, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	value, err := c.eventService.Get(eventId)
	if err != nil {
		return nil, err
	}
	if value == nil || len(value.Guests) == 0 {
		return []*app.Guest{}, nil
	}
	return value.Guests, nil
}

func (c *GuestService) CreateOrUpdate(eventId string, u *app.Guest) (*app.Guest, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := u.Validate()
	if err != nil {
		return nil, err
	}
	// If Id is nil populate it
	if u.Id == "" {
		u.Id = app.GenerateId(u.FirstName + u.LastName)
	}

	event, err := c.eventService.Get(eventId)
	if err != nil {
		return nil, err
	}
	// Now try to get the event
	_, idx, err := searchByGuestId(event.Guests, u.Id)
	if err != nil {
		u.TimeCreatedOn = time.Now()
		u.TimeUpdatedOn = time.Now()
		event.Guests = append(event.Guests, u)
	} else {
		// If it exists update the time stamp and return, we should be more strict about validations but dah!
		event.Guests[idx] = u
		event.Guests[idx].TimeUpdatedOn = time.Now()
	}

	log.Printf("Created event %s", u.Id)
	return u, nil
}

func (c *GuestService) Delete(eventId, id string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	event, _ := c.eventService.Get(eventId)

	_, i, err := searchByGuestId(event.Guests, id)
	if err != nil {
		return errors.New("object event does not exist")
	}
	event.Guests = slices.Delete(event.Guests, i, i+1)
	return nil
}

func searchByGuestId(guests []*app.Guest, id string) (*app.Guest, int, error) {
	for i, value := range guests {
		if value.Id == id {
			return value, i, nil
		}
	}
	return nil, -1, errors.New("guest not found")
}
