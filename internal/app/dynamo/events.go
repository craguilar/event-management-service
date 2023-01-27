package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
