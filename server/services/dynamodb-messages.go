package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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

func CreateMessage(client *dynamodb.Client, tableName string, msg Message) error {
	item, err := attributevalue.MarshalMap(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(id)"), // prevent overwrite
	})
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}

func GetChatMessage(client *dynamodb.Client, tableName, chatID, msgID string) (*Message, error) {
	out, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"chatId": &types.AttributeValueMemberS{Value: chatID},
			"id":     &types.AttributeValueMemberS{Value: msgID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	if out.Item == nil {
		return nil, fmt.Errorf("message not found")
	}

	var msg Message
	err = attributevalue.UnmarshalMap(out.Item, &msg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}
	return &msg, nil
}

func GetAllChatMessages(client *dynamodb.Client, tableName, chatID string) ([]Message, error) {
	out, err := client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("chatId = :c"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":c": &types.AttributeValueMemberS{Value: chatID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}

	var messages []Message
	err = attributevalue.UnmarshalListOfMaps(out.Items, &messages)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal messages: %w", err)
	}
	return messages, nil
}

func UpdateMessage(client *dynamodb.Client, tableName, chatID, msgID string, updates map[string]types.AttributeValue) error {
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
			"chatId": &types.AttributeValueMemberS{Value: chatID},
			"id":     &types.AttributeValueMemberS{Value: msgID},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprAttrValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}
	return nil
}

func DeleteMessage(client *dynamodb.Client, tableName, chatID, msgID string) error {
	_, err := client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"chatId": &types.AttributeValueMemberS{Value: chatID},
			"id":     &types.AttributeValueMemberS{Value: msgID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}
