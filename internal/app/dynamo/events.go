package dynamo

import (
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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

func InitLocalDb(overrideUrl, tableName string) *DBConfig {

	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String("dummy"),
		Endpoint:    aws.String(overrideUrl),
		Credentials: credentials.NewStaticCredentials("AKID", "SECRET_KEY", "TOKEN"),
	})
	if err != nil {
		log.Fatalf("Error found %s", err)
	}
	return InitDb(dynamodb.New(awsSession), tableName)
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
			"entityType": {
				S: aws.String("EVENT-" + id),
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
		IndexName: aws.String("ownerIdx"),
		KeyConditions: map[string]*dynamodb.Condition{
			"entityType": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String("OWNER-" + userName),
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
	aEvent["id"] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	aEvent["entityType"] = &dynamodb.AttributeValue{S: aws.String("EVENT-" + u.Id)}

	aOwner, err := dynamodbattribute.MarshalMap(eventOwner(eventManager, u))
	if err != nil {
		return nil, err
	}
	aOwner["id"] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	aOwner["entityType"] = &dynamodb.AttributeValue{S: aws.String("OWNER-" + eventManager)}

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
	return errors.New("not implemented")
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
