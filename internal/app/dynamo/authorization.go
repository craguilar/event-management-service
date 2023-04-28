package dynamo

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/craguilar/event-management-service/internal/app"
)

// EventService represents a Dynamo DB implementation of internal.EventService.
type AuthorizationService struct {
	db *DBConfig
}

func NewAuthorizationService(db *DBConfig) *AuthorizationService {
	if db == nil {
		log.Panicf("Null reference to db config in EventService")
	}
	return &AuthorizationService{
		db: db,
	}
}

func (a *AuthorizationService) Authorize(userName, eventId string) bool {
	userName = strings.ToUpper(userName)
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			a.db.PK_ID: {
				S: aws.String(eventId),
			},
			a.db.SORT_KEY: {
				S: aws.String(_SORT_KEY_OWNER_PREFIX + userName),
			},
		},
		TableName: &a.db.TableName,
	}

	result, err := a.db.DbService.GetItem(input)
	if err != nil {
		log.Printf("Error when GetItem for authorize %s", err)
		return false

	}
	owner := &app.EventOwner{}
	err = dynamodbattribute.UnmarshalMap(result.Item, owner)
	if err != nil {
		log.Printf("Error when UnmarshalMap for authorize in eventService %s", err)
		return false
	}
	if owner.EventSummary.Id == "" {
		log.Printf("Error when UnmarshalMap for authorize in eventService %s", err)
		return false
	}
	return true
}
