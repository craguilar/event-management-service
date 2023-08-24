package dynamo

import (
	"log"
	"time"

	"github.com/craguilar/event-management-service/internal/app"
)

// EventService represents a Dynamo DB implementation of internal.EventService.
type EventActions struct {
	db                  *DBConfig
	eventService        *EventService
	taskService         *TaskService
	notificationService *app.EmailNotificationService
}

func NewEventActionsService(db *DBConfig, event *EventService, task *TaskService, notification *app.EmailNotificationService) *EventActions {
	if db == nil {
		log.Panicf("Null reference to db config in EventService")
	}
	return &EventActions{
		db:                  db,
		eventService:        event,
		taskService:         task,
		notificationService: notification,
	}
}

func (c *EventActions) SendPendingTasksNotifications() error {
	events, err := c.eventService.ListBy(func(event *app.EventSummary) bool {
		return event.EventDay.After(time.Now()) && event.NotificationEnabled
	})
	if err != nil {
		return err
	}
	for _, event := range events {

		// Get the tasks
		tasks, err := c.taskService.List(event.Id)
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			continue
		}
		template := app.TemplatePendingTasksNotifications(event.Name, tasks)
		// Get owners
		owners, err := c.eventService.ListOwners(event.Id)
		if err != nil {
			return err
		}
		log.Printf("Sending notification for %s to %d recipients with %d tasks", event.Name, len(owners.SharedEmails), len(tasks))
		for _, toEmail := range owners.SharedEmails {
			err = c.notificationService.SendEmailNotification(toEmail, "Pending Tasks for "+event.Name, template.String())
			if err != nil {
				log.Printf("WARN: Failed to send notification for %s with error %s", toEmail, err)
			}
		}

	}
	return nil
}

func TemplatePendingTasksNotifications(s string, tasks []*app.Task) {
	panic("unimplemented")
}
