package dynamo

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/craguilar/event-management-service/internal/app"
)

const _SORT_KEY_EVENT_PREFIX = "EVENT-"
const _SORT_KEY_OWNER_PREFIX = "OWNER-"

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
			c.db.PK_ID: {
				S: aws.String(id),
			},
			c.db.SORT_KEY: {
				S: aws.String(_SORT_KEY_EVENT_PREFIX + id),
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
	log.Printf("Return %v", event)
	return event, nil
}

func (c *EventService) List(userName string) ([]*app.EventSummary, error) {

	log.Printf("Getting all events for %s", userName)
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(c.db.TableName),
		IndexName: aws.String(c.db.GSI_OWNER),
		KeyConditions: map[string]*dynamodb.Condition{
			c.db.SORT_KEY: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(_SORT_KEY_OWNER_PREFIX + userName),
					},
				},
			},
		},
	}

	result, err := c.db.DbService.Query(queryInput)
	if err != nil {
		return nil, err
	}
	items := new([]app.EventOwner)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, err
	}
	list := []*app.EventSummary{}
	for _, value := range *items {
		list = append(list, &app.EventSummary{Id: value.EventSummary.Id, Name: value.EventSummary.Name, MainLocation: value.EventSummary.MainLocation, EventDay: value.EventSummary.EventDay, TimeCreatedOn: value.EventSummary.TimeCreatedOn})
	}
	return list, nil
}

func (c *EventService) CreateOrUpdate(eventManager string, u *app.Event) (*app.Event, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	// If Id is nil populate it
	if u.Id == "" {
		u.Id = app.GenerateId(u.Name)
	}

	log.Printf("CreateOrUpdate event with name %s /%s", u.Name, u.Id)

	value, err := c.Get(u.Id)
	if err != nil {
		return nil, err
	}
	if value == nil {
		u.TimeCreatedOn = time.Now()
	}
	// If it exists update the time stamp!
	u.TimeUpdatedOn = time.Now()
	aEvent, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	// Assign dynamo db key
	aEvent[c.db.PK_ID] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	aEvent[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(_SORT_KEY_EVENT_PREFIX + u.Id)}

	aOwner, err := dynamodbattribute.MarshalMap(eventOwner(eventManager, u))
	if err != nil {
		return nil, err
	}
	aOwner[c.db.PK_ID] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	aOwner[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(_SORT_KEY_OWNER_PREFIX + eventManager)}

	transactions := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				Put: &dynamodb.Put{
					Item:      aEvent,
					TableName: &c.db.TableName,
				},
			},
			{
				Put: &dynamodb.Put{
					Item:      aOwner,
					TableName: &c.db.TableName,
				},
			},
		},
	}
	output, err := c.db.DbService.TransactWriteItems(transactions)

	if err != nil {
		return nil, err
	}

	log.Printf("Created event with name %s /%s - output %s", u.Name, u.Id, output)
	return u, nil
}

func (c *EventService) Delete(id string) error {

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			c.db.PK_ID: {
				S: aws.String(id),
			},
		},
		TableName: &c.db.TableName,
	}

	_, err := c.db.DbService.DeleteItem(input)
	if err != nil {
		log.Printf("Got error calling DeetItem:")
		return err
	}
	return nil
}

func eventOwner(userName string, event *app.Event) *app.EventOwner {

	return &app.EventOwner{
		OwnerEmail: userName,
		EventSummary: &app.EventSummary{
			Id:            event.Id,
			Name:          event.Name,
			MainLocation:  event.MainLocation,
			EventDay:      event.EventDay,
			TimeCreatedOn: event.TimeCreatedOn,
		},
	}
}
