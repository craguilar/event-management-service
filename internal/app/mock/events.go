package mock

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/craguilar/event-management-service/internal/app"
)

type EventService struct {
	db   map[string]*app.Event
	lock sync.RWMutex
}

func NewEventService() *EventService {
	return &EventService{
		db: make(map[string]*app.Event),
	}
}

func (c *EventService) Get(id string) (*app.Event, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	// Get the value
	value, exists := c.db[id]
	if !exists {
		return nil, nil
	}
	return value, nil
}

func (c *EventService) List(user string) ([]*app.EventSummary, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	list := []*app.EventSummary{}
	for _, value := range c.db {
		list = append(list, &app.EventSummary{Id: value.Id, Name: value.Name, MainLocation: value.MainLocation, EventDay: value.EventDay, TimeCreatedOn: value.TimeCreatedOn})
	}
	return list, nil
}

func (c *EventService) CreateOrUpdate(u *app.Event) (*app.Event, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := u.Validate()
	if err != nil {
		return nil, err
	}
	// If Id is nil populate it
	if u.Id == "" {
		data := []byte(strings.ToUpper(u.Name))
		u.Id = fmt.Sprintf("%x", md5.Sum(data))
		if err != nil {
			return nil, err
		}
	}
	_, exists := c.db[u.Id]
	if !exists {
		u.TimeCreatedOn = time.Now()
		u.TimeUpdatedOn = time.Now()
		c.db[u.Id] = u
		return u, nil
	}
	// If it exists update the time stamp and return, we should be more strict about validations but dah!
	u.TimeUpdatedOn = time.Now()
	c.db[u.Id] = u
	log.Printf("Created event %s", u.Id)
	return u, nil
}

func (c *EventService) Delete(id string) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, exists := c.db[id]
	if !exists {
		return errors.New("object event does not exist")
	}
	delete(c.db, id)
	return nil
}