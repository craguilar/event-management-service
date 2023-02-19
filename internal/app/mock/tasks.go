package mock

import (
	"errors"

	"github.com/craguilar/event-management-service/internal/app"
)

type TaskService struct {
}

func (c *TaskService) Get(eventId, id string) (*app.Task, error) {
	return nil, errors.New("not implemented")
}

func (c *TaskService) List(eventId string) ([]*app.Task, error) {
	return nil, errors.New("not implemented")
}

func (c *TaskService) CreateOrUpdate(eventId string, u *app.Task) (*app.Task, error) {
	return nil, errors.New("not implemented")
}

func (c *TaskService) Delete(eventId, id string) error {
	return errors.New("not implemented")
}
