package dynamo

import (
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/craguilar/event-management-service/internal/app"
)

const _SORT_KEY_GUEST_PREFIX = "GUEST-"

type GuestService struct {
	db *DBConfig
}

func NewGuestService(db *DBConfig) *GuestService {
	return &GuestService{
		db: db,
	}
}

func (c *GuestService) Get(eventId, id string) (*app.Guest, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(eventId),
			},
			"entityType": {
				S: aws.String(id),
			},
		},
		TableName: &c.db.TableName,
	}

	result, err := c.db.DbService.GetItem(input)
	if err != nil {
		return nil, err

	}

	event := &app.Guest{}
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

func (c *GuestService) List(eventId string) ([]*app.Guest, error) {
	log.Printf("Getting all events for %s", eventId)
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(c.db.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			c.db.PK_ID: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(eventId),
					},
				},
			},
			c.db.SORT_KEY: {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(_SORT_KEY_GUEST_PREFIX),
					},
				},
			},
		},
	}

	result, err := c.db.DbService.Query(queryInput)
	if err != nil {
		return nil, err
	}
	// Given we use a single table model until here result.Items contains ALL the same duplicated Id :)
	//err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	list := []*app.Guest{}
	for _, value := range result.Items {
		guest := &app.Guest{}
		err = dynamodbattribute.UnmarshalMap(value, guest)
		if err != nil {
			return nil, err
		}
		guest.Id = *aws.String(*value[c.db.SORT_KEY].S)
		list = append(list, guest)
	}
	return list, nil
}

func (c *GuestService) CreateOrUpdate(eventId string, u *app.Guest) (*app.Guest, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	// If Id is nil populate it
	if u.Id == "" {
		u.Id = app.GenerateId(strings.ToUpper(u.FirstName + u.LastName))
	}

	log.Printf("CreateOrUpdate guest with Id /%s", u.Id)

	value, err := c.Get(eventId, u.Id)
	if err != nil {
		return nil, err
	}
	if value == nil {
		u.TimeCreatedOn = time.Now()
	}
	// If it exists update the time stamp!
	u.TimeUpdatedOn = time.Now()
	aGuest, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	// Assign dynamo db key
	aGuest[c.db.PK_ID] = &dynamodb.AttributeValue{S: aws.String(eventId)}
	aGuest[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(_SORT_KEY_GUEST_PREFIX + u.Id)}
	input := &dynamodb.PutItemInput{
		Item:      aGuest,
		TableName: &c.db.TableName,
	}

	_, err = c.db.DbService.PutItem(input)
	if err != nil {
		return nil, err
	}
	log.Printf("Creatde guest with Id %s", u.Id)
	return u, nil
}

func (c *GuestService) Delete(eventId, id string) error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			c.db.PK_ID: {
				S: aws.String(eventId),
			},
			c.db.SORT_KEY: {
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
