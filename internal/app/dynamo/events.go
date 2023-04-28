package dynamo

import (
	"errors"
	"log"
	"strings"
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
	db        *DBConfig
	authorize *AuthorizationService
}

func NewEventService(db *DBConfig, authorize *AuthorizationService) *EventService {
	if db == nil {
		log.Panicf("Null reference to db config in EventService")
	}
	return &EventService{
		db:        db,
		authorize: authorize,
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
	userName = strings.ToUpper(userName)
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

// Added AuthZ to prevent the situation where, someone hacks its way and sends the
// eventId same as an existing eventId from other owner , in the current situation we will accept
// it AND end up adding a new owner . Ideally before calling CreateOrUpdate we should check if
// updating an existing event , the eventManager MUST match an existing OWNER in the table.
func (c *EventService) CreateOrUpdate(eventManager string, u *app.Event) (*app.Event, error) {
	eventManager = strings.ToUpper(eventManager)
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	if u.Id != "" && !c.authorize.Authorize(eventManager, u.Id) {
		return nil, errors.New("unauthorized")
	}
	// TODO: Document why I decided to add a random Id
	if u.Id == "" {
		u.Id, err = app.GenerateRandomId()
		if err != nil {
			return nil, err
		}
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

func (c *EventService) Delete(eventManager, id string) error {

	if !c.authorize.Authorize(eventManager, id) {
		return errors.New("unauthorized")
	}
	// Get ALL associated elements
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(c.db.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			c.db.PK_ID: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(id),
					},
				},
			},
		},
	}

	result, err := c.db.DbService.Query(queryInput)
	if err != nil {
		log.Printf("Error when querying by HASH key - %s", err)
		return err
	}

	transactions := []*dynamodb.TransactWriteItem{}
	for _, value := range result.Items {
		transactions = append(transactions, &dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				Key: map[string]*dynamodb.AttributeValue{
					c.db.PK_ID: {
						S: aws.String(id),
					},
					c.db.SORT_KEY: {
						S: aws.String(*value[c.db.SORT_KEY].S),
					},
				},
				TableName: &c.db.TableName,
			},
		})

	}

	// Batch delete
	transactWriteInput := &dynamodb.TransactWriteItemsInput{TransactItems: transactions}
	_, err = c.db.DbService.TransactWriteItems(transactWriteInput)
	if err != nil {
		log.Printf("Got error calling Delete event - %s", err)
		return err
	}
	return nil
}

// Owner
func (c *EventService) ListOwners(id string) (*app.EventSharedEmails, error) {

	log.Printf("Getting all events for %s", id)
	var queryInput = &dynamodb.QueryInput{
		TableName: aws.String(c.db.TableName),
		KeyConditions: map[string]*dynamodb.Condition{
			c.db.PK_ID: {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(id),
					},
				},
			},
			c.db.SORT_KEY: {
				ComparisonOperator: aws.String("BEGINS_WITH"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(_SORT_KEY_OWNER_PREFIX),
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
	sharedEmails := &app.EventSharedEmails{
		EventId:      id,
		SharedEmails: []string{},
	}
	for _, value := range result.Items {
		sortKey := *aws.String(*value[c.db.SORT_KEY].S)
		if !strings.HasPrefix(sortKey, _SORT_KEY_OWNER_PREFIX) {
			continue
		}
		sharedEmails.SharedEmails = append(sharedEmails.SharedEmails, strings.ReplaceAll(sortKey, _SORT_KEY_OWNER_PREFIX, ""))
	}
	return sharedEmails, nil
}

// AddOwner receives a current eventManager coming from Authorization token AND adds
// a new OWNER to an eventId.
func (c *EventService) CreateOwner(userName string, u *app.EventSharedEmails) (*app.EventSharedEmails, error) {
	// Does eventManager check
	if !c.authorize.Authorize(userName, u.EventId) {
		return nil, errors.New("unauthorized")
	}
	event, err := c.Get(u.EventId)
	if err != nil {
		return nil, err
	}
	//
	transactions := []*dynamodb.TransactWriteItem{}
	for i := 0; i < len(u.SharedEmails); i++ {
		u.SharedEmails[i] = strings.ToUpper(u.SharedEmails[i])
		if strings.HasPrefix(u.SharedEmails[i], _SORT_KEY_OWNER_PREFIX) {
			return nil, errors.New("plain emails expected")
		}
		aOwner, err := dynamodbattribute.MarshalMap(eventOwner(u.SharedEmails[i], event))
		if err != nil {
			return nil, err
		}
		aOwner[c.db.PK_ID] = &dynamodb.AttributeValue{S: aws.String(u.EventId)}
		aOwner[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(_SORT_KEY_OWNER_PREFIX + u.SharedEmails[i])}

		transactions = append(transactions, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:      aOwner,
				TableName: &c.db.TableName,
			},
		})
	}

	if err != nil {
		return nil, err
	}
	transactWriteInput := &dynamodb.TransactWriteItemsInput{TransactItems: transactions}
	_, err = c.db.DbService.TransactWriteItems(transactWriteInput)
	if err != nil {
		log.Printf("Got error calling Delete - %s", err)
		return nil, err
	}
	return u, nil
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
