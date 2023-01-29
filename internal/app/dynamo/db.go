package dynamo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
		PK_ID:     "id",
		SORT_KEY:  "entityType",
		GSI_OWNER: "ownerIdx",
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
