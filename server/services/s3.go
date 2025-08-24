package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func ConnectS3() *s3.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	s3_region := os.Getenv("AWS_REGION_S3")

	s3Cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(s3_region),
		config.WithCredentialsProvider(
			aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)
	s3Client := s3.NewFromConfig(s3Cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	_, err = s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("unable to load S3 buckets, %v", err)
	}
	log.Printf("Connected to S3\n")
	return s3Client
}

func UploadFile(client *s3.Client, filename string, fileContent multipart.File) (string, error) {
	bucketName := os.Getenv("AWS_BUCKET")
	s3_region := os.Getenv("AWS_REGION_S3")

	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("uploads/" + filename),
		Body:   fileContent,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Use regional endpoint for correct URL
	fileURL := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/uploads/%s", bucketName, s3_region, filename)
	return fileURL, nil
}

func DownloadFile(client *s3.Client, filename string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	bucketName := os.Getenv("AWS_BUCKET")

	fileKey := filename

	expiration := time.Duration(5) * time.Minute

	presignClient := s3.NewPresignClient(client)
	presignedURL, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("failed to presign URL: %w", err)
	}

	return presignedURL.URL, nil
}
