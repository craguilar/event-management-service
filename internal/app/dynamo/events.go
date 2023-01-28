package dynamo

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/craguilar/event-management-service/internal/app"
)

type DBConfig struct {
	DbService *dynamodb.DynamoDB
	TableName string
}

func InitDb(db *dynamodb.DynamoDB, tableName string) *DBConfig {
	return &DBConfig{
		DbService: db,
		TableName: tableName,
	}
}

// EventService represents a Dynamo DB implementation of internal.EventService.
type EventService struct {
	db *DBConfig
}

func NewEventService(db *DBConfig) *EventService {
	return &EventService{
		db: db,
	}
}

func (c *EventService) Get(id string) (*app.Event, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: &c.db.TableName,
	}

	result, err := c.db.DbService.GetItem(input)
	if err != nil {
		return nil, err

	}

	event := &app.Event{}
	err = dynamodbattribute.UnmarshalMap(result.Item, event)
	if err != nil {
		return nil, err
	}
	if event.Id == "" {
		return nil, nil
	}
	log.Printf("Return car %v", event)
	return event, nil
}

func (c *EventService) List(user string) ([]*app.EventSummary, error) {
	input := &dynamodb.ScanInput{
		TableName: &c.db.TableName,
	}
	result, err := c.db.DbService.Scan(input)
	if err != nil {
		return nil, err
	}
	items := new([]app.Event)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, err
	}
	list := []*app.EventSummary{}
	for _, value := range *items {
		list = append(list, &app.EventSummary{Id: value.Id, Name: value.Name, MainLocation: value.MainLocation, EventDay: value.EventDay, TimeCreatedOn: value.TimeCreatedOn})
	}
	return list, nil
}

func (c *EventService) CreateOrUpdate(u *app.Event) (*app.Event, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	log.Printf("Creating car with plate %s", u.Id)
	// TODO: Check if user exists
	value, err := c.Get(u.Id)
	if err != nil {
		return nil, err
	}
	if value == nil {
		u.TimeCreatedOn = time.Now()
	}
	// If it exists update the time stamp!
	u.TimeUpdatedOn = time.Now()
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	// Assign dynamo db key
	av["id"] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &c.db.TableName,
	}

	output, err := c.db.DbService.PutItem(input)
	if err != nil {
		return nil, err
	}
	log.Printf("Created %s", output)
	return u, nil
}

func (c *EventService) Delete(id string) error {
	return errors.New("not implemented")
}
