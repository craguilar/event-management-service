package dynamo

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/craguilar/event-management-service/internal/app"
)

const _SORT_KEY_TASK_PREFIX = "TASK-"

type TaskService struct {
	db *DBConfig
}

func NewTaskService(db *DBConfig) *TaskService {
	return &TaskService{
		db: db,
	}
}

func (c *TaskService) Get(eventId, id string) (*app.Task, error) {
	input := &dynamodb.GetItemInput{
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

	result, err := c.db.DbService.GetItem(input)
	if err != nil {
		return nil, err

	}

	task := &app.Task{}
	err = dynamodbattribute.UnmarshalMap(result.Item, task)
	if err != nil {
		return nil, err
	}
	if task.Id == "" {
		return nil, nil
	}
	log.Printf("Return %v", task)
	return task, nil
}

func (c *TaskService) List(eventId string) ([]*app.Task, error) {
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
						S: aws.String(_SORT_KEY_TASK_PREFIX),
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
	list := []*app.Task{}
	for _, value := range result.Items {
		task := &app.Task{}
		err = dynamodbattribute.UnmarshalMap(value, task)
		if err != nil {
			return nil, err
		}
		task.Id = *aws.String(*value[c.db.SORT_KEY].S)
		list = append(list, task)
	}
	return list, nil
}

func (c *TaskService) CreateOrUpdate(eventId string, u *app.Task) (*app.Task, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	// If Id is nil populate it
	if u.Id == "" {
		u.Id, err = app.GenerateRandomId()
		if err != nil {
			return nil, err
		}
	}

	log.Printf("CreateOrUpdate guest with Id /%s", u.Id)
	aTask, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	value, err := c.Get(eventId, u.Id)
	if err != nil {
		return nil, err
	}
	// Create else Update
	aTask[c.db.PK_ID] = &dynamodb.AttributeValue{S: aws.String(eventId)}
	if value == nil {
		u.TimeCreatedOn = time.Now()
		aTask[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(_SORT_KEY_TASK_PREFIX + u.Id)}
	} else {
		u.TimeUpdatedOn = time.Now()
		aTask[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	}
	input := &dynamodb.PutItemInput{
		Item:      aTask,
		TableName: &c.db.TableName,
	}

	_, err = c.db.DbService.PutItem(input)
	if err != nil {
		return nil, err
	}
	log.Printf("Creatde guest with Id %s", u.Id)
	return u, nil
}

func (c *TaskService) Delete(eventId, id string) error {
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
		log.Printf("Got error calling delete task  %s ", err)
		return err
	}
	return nil
}
