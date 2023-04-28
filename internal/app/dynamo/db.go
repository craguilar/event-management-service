package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const C_PK_ID = "id"
const C_SORT_KEY = "entityType"
const C_GSI_OWNER = "ownerIdx"

type DBConfig struct {
	DbService *dynamodb.DynamoDB
	TableName string
	PK_ID     string
	SORT_KEY  string
	GSI_OWNER string
}

func InitDb(db *dynamodb.DynamoDB, tableName string) *DBConfig {
	return &DBConfig{
		DbService: db,
		TableName: tableName,
		PK_ID:     C_PK_ID,
		SORT_KEY:  C_SORT_KEY,
		GSI_OWNER: C_GSI_OWNER,
	}
}

// DO NOT USE IN PRODUCTION
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
