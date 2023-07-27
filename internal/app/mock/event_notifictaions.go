package mock

import (
	"log"
	"time"

	"github.com/craguilar/event-management-service/internal/app"
)

// EventService represents a Dynamo DB implementation of internal.EventService.
type EventActions struct {
	eventService *EventService
	taskService  *TaskService
}

func NewEventActionsService(event *EventService, task *TaskService) *EventActions {

	return &EventActions{
		eventService: event,
		taskService:  task,
	}
}

func (c *EventActions) SendPendingTasksNotifications() error {
	events, err := c.eventService.ListBy(func(event *app.EventSummary) bool {
		return event.EventDay.After(time.Now())
	})
	if err != nil {
		return err
	}
	for _, event := range events {

		tasks, err := c.taskService.List(event.Id)
		if err != nil {
			return err
		}
		for _, task := range tasks {
			log.Printf("%v", task)
			// Do something
		}
	}
	return nil
}
