package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func CreateChatsTable(client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash, // Partition Key
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	return err
}

func CreateChat(client *dynamodb.Client, tableName string, chat Chat) error {
	item, err := attributevalue.MarshalMap(chat)
	if err != nil {
		return fmt.Errorf("failed to marshal chat: %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"), // prevent overwrite
	})
	if err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}
	return nil
}

func GetAllChats(client *dynamodb.Client, tableName string) ([]Chat, error) {
	out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan chats: %w", err)
	}

	var chats []Chat
	err = attributevalue.UnmarshalListOfMaps(out.Items, &chats)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chats: %w", err)
	}
	return chats, nil
}

func GetChatById(client *dynamodb.Client, tableName, chatID string) (*Chat, error) {
	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: chatID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}
	if out.Item == nil {
		return nil, fmt.Errorf("chat not found")
	}

	var chat Chat
	err = attributevalue.UnmarshalMap(out.Item, &chat)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal chat: %w", err)
	}
	return &chat, nil
}

func UpdateChat(client *dynamodb.Client, tableName, chatID string, updates map[string]types.AttributeValue) error {
	updateExpr := "SET"
	exprAttrValues := map[string]types.AttributeValue{}
	i := 0

	for k, v := range updates {
		if i > 0 {
			updateExpr += ","
		}
		placeholder := fmt.Sprintf(":%s", k)
		updateExpr += fmt.Sprintf(" %s = %s", k, placeholder)
		exprAttrValues[placeholder] = v
		i++
	}

	_, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: chatID},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return fmt.Errorf("failed to update chat: %w", err)
	}
	return nil
}

func DeleteChat(client *dynamodb.Client, tableName, chatID string) error {
	_, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: chatID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete chat: %w", err)
	}
	return nil
}
