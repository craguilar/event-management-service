package mock

import (
	"testing"
	"time"

	"github.com/craguilar/event-management-service/internal/app"
)

func TestList(t *testing.T) {
	eventService := NewEventService()

	event := &app.Event{Name: "My Birthday", MainLocation: "Golden Gate Park", EventDay: time.Now()}
	event, err := eventService.CreateOrUpdate(event)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	guestService := NewGuestService(eventService)
	//Create
	guest := &app.Guest{FirstName: "Mickey", LastName: "Mouse", Tentative: false, NumberOfSeats: 1}
	_, err = guestService.CreateOrUpdate(event.Id, guest)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}

	guests, err := guestService.List(event.Id)
	// Assertion

	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	if len(guests) != 1 {
		t.Fatalf("A list guest of size 1 is expected")
	}
}

func TestGuestCreation(t *testing.T) {
	eventService := NewEventService()

	event := &app.Event{Name: "My Birthday", MainLocation: "Golden Gate Park", EventDay: time.Now()}
	event, err := eventService.CreateOrUpdate(event)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	guestService := NewGuestService(eventService)

	guest := &app.Guest{FirstName: "Mickey", LastName: "Mouse", Tentative: false, NumberOfSeats: 1}
	_, err = guestService.CreateOrUpdate(event.Id, guest)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	event2, err := eventService.Get(event.Id)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	// Assertion
	if len(event2.Guests) == 0 {
		t.Fatalf("A new guest expectde nothing found %s", err)
	}
}

func TestGuestUpdate(t *testing.T) {
	eventService := NewEventService()

	event := &app.Event{Name: "My Birthday", MainLocation: "Golden Gate Park", EventDay: time.Now()}
	event, err := eventService.CreateOrUpdate(event)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	guestService := NewGuestService(eventService)
	//Create
	guest := &app.Guest{FirstName: "Mickey", LastName: "Mouse", Tentative: false, NumberOfSeats: 1}
	_, err = guestService.CreateOrUpdate(event.Id, guest)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	//Update
	guest.Email = "my-email@nowhere.com"
	guestService.CreateOrUpdate(event.Id, guest)
	// Assertion
	event2, err := eventService.Get(event.Id)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	if len(event2.Guests) == 0 || event2.Guests[0].Email != "my-email@nowhere.com" {
		t.Fatalf("A new guest expectde nothing found %s", err)
	}
}

func TestDelete(t *testing.T) {
	eventService := NewEventService()

	event := &app.Event{Name: "My Birthday", MainLocation: "Golden Gate Park", EventDay: time.Now()}
	event, err := eventService.CreateOrUpdate(event)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	guestService := NewGuestService(eventService)
	//Create
	guest := &app.Guest{FirstName: "Mickey", LastName: "Mouse", Tentative: false, NumberOfSeats: 1}
	_, err = guestService.CreateOrUpdate(event.Id, guest)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}

	guests, err := guestService.List(event.Id)
	// Then pre condition Assertions
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	if len(guests) != 1 {
		t.Fatalf("A list guest of size 1 is expected")
	}
	guestService.Delete(event.Id, guests[0].Id)
	//Assert
	guests2, err := guestService.List(event.Id)
	if err != nil {
		t.Fatalf("Test failed with error %s", err)
	}
	if len(guests2) != 0 {
		t.Fatalf("A list guest of size 0 is expected")
	}
}
