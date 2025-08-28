package services

type User struct {
	ID       string `json:"id" dynamodbav:"id"`
	Name     string `json:"name" dynamodbav:"name"`
	Email    string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
}

type Chat struct {
	ID          string   `json:"id" dynamodbav:"id"`
	Users       []string `json:"users" dynamodbav:"users"`
	Messages    []string `json:"messages" dynamodbav:"messages"`
	DateCreated int64    `json:"dateCreated" dynamodbav:"dateCreated"`
	DateUpdated int64    `json:"dateUpdated" dynamodbav:"dateUpdated"`
}

type Message struct {
	ID        string   `json:"id" dynamodbav:"id"`
	ChatID    string   `json:"chatId" dynamodbav:"chatId"`
	SenderID  string   `json:"senderId" dynamodbav:"senderId"`
	Content   string   `json:"content" dynamodbav:"content"`
	Media     []string `json:"media" dynamodbav:"media"`
	Timestamp int64    `json:"timestamp" dynamodbav:"timestamp"`
}

type UserFile struct {
	UserID   string `dynamodbav:"userId"` // partition key
	FileID   string `dynamodbav:"fileId"` // sort key
	FileKey  string `dynamodbav:"fileKey"`
	Uploaded int64  `dynamodbav:"uploaded"`
}
