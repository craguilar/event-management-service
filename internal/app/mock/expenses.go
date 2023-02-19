package mock

import (
	"errors"

	"github.com/craguilar/event-management-service/internal/app"
)

type ExpenseService struct {
}

func (c *ExpenseService) Get(eventId, id string) (*app.ExpenseCategory, error) {
	return nil, errors.New("not implemented")
}

func (c *ExpenseService) List(eventId string) ([]*app.ExpenseCategory, error) {
	return nil, errors.New("not implemented")

}
func (c *ExpenseService) CreateOrUpdate(eventId string, u *app.ExpenseCategory) (*app.ExpenseCategory, error) {

	return nil, errors.New("not implemented")
}
func (c *ExpenseService) Delete(eventId, id string) error {
	return errors.New("not implemented")

}
