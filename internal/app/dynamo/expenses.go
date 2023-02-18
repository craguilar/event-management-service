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

const _SORT_KEY_EXPENSE_CATEGORY_PREFIX = "EXPENSE_CATEGORY-"

type ExpenseService struct {
	db *DBConfig
}

func NewExpenseService(db *DBConfig) *ExpenseService {
	return &ExpenseService{
		db: db,
	}
}

func (c *ExpenseService) Get(eventId, id string) (*app.ExpenseCategory, error) {
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

	category := &app.ExpenseCategory{}
	err = dynamodbattribute.UnmarshalMap(result.Item, category)
	if err != nil {
		return nil, err
	}
	if category.Category == "" {
		return nil, nil
	}
	log.Printf("Return %v", category)
	return category, nil
}

func (c *ExpenseService) List(eventId string) ([]*app.ExpenseCategory, error) {
	log.Printf("Getting all expenses for %s", eventId)
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
						S: aws.String(_SORT_KEY_EXPENSE_CATEGORY_PREFIX),
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
	list := []*app.ExpenseCategory{}
	for _, value := range result.Items {
		expense := &app.ExpenseCategory{}
		err = dynamodbattribute.UnmarshalMap(value, expense)
		if err != nil {
			return nil, err
		}
		expense.Id = *aws.String(*value[c.db.SORT_KEY].S)
		list = append(list, expense)
	}
	return list, nil
}

// TODO: We shall add a condition that restricts the max number of
func (c *ExpenseService) CreateOrUpdate(eventId string, u *app.ExpenseCategory) (*app.ExpenseCategory, error) {
	err := u.Validate()
	if err != nil {
		return nil, err
	}
	// If Id is nil populate it
	if u.Id == "" {
		u.Id = app.GenerateId(strings.ToUpper(u.Category))
	}
	amountPaid := 0.0
	for i := range u.Expenses {
		amountPaid += u.Expenses[i].AmountPaid
		if u.Expenses[i].Id == "" {
			u.Expenses[i].Id, err = app.GenerateRandomId()
			if err != nil {
				return nil, err
			}
		}
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
	u.AmountPaid = amountPaid
	aExpense, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	// Assign dynamo db key
	aExpense[c.db.PK_ID] = &dynamodb.AttributeValue{S: aws.String(eventId)}
	if value == nil {
		u.TimeCreatedOn = time.Now()
		u.Expenses = []*app.Expense{}
		aExpense[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(_SORT_KEY_EXPENSE_CATEGORY_PREFIX + u.Id)}
	} else {
		u.TimeUpdatedOn = time.Now()
		aExpense[c.db.SORT_KEY] = &dynamodb.AttributeValue{S: aws.String(u.Id)}
	}
	input := &dynamodb.PutItemInput{
		Item:      aExpense,
		TableName: &c.db.TableName,
	}
	_, err = c.db.DbService.PutItem(input)
	if err != nil {
		return nil, err
	}
	log.Printf("Creatde an expense with Id %s", u.Id)
	return u, nil
}

func (c *ExpenseService) Delete(eventId, id string) error {
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
