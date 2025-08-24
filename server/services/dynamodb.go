package services

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

func ConnectDB() *dynamodb.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	ddb_region := os.Getenv("AWS_REGION_DDB")

	ddbCfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ddb_region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		),
	)
	ddbClient := dynamodb.NewFromConfig(ddbCfg)
	_, err = GetTables(ddbClient)
	if err != nil {
		log.Fatalf("Error connecting to DynamoDB.")
	}
	log.Printf("Connected to DynamoDB\n")
	return ddbClient
}

func GetTables(client *dynamodb.Client) ([]string, error) {

	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return nil, err
	}
	return result.TableNames, nil
}
