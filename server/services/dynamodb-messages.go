package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateMessagesTable(client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("chatId"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("chatId"),
				KeyType:       types.KeyTypeHash, // Partition Key
			},
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeRange, // Sort Key
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	return err
}
