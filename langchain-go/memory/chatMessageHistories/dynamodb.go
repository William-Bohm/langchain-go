package chatMessageHistories

import (
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBChatMessageHistory struct {
	table     *dynamodb.DynamoDB
	sessionID string
	tableName string
}

func NewDynamoDBChatMessageHistory(tableName, sessionID string) *DynamoDBChatMessageHistory {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	dynamoDB := dynamodb.New(sess)

	return &DynamoDBChatMessageHistory{
		table:     dynamoDB,
		sessionID: sessionID,
		tableName: tableName,
	}
}

func (d *DynamoDBChatMessageHistory) Messages() ([]rootSchema.BaseMessageInterface, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"SessionId": {
				S: aws.String(d.sessionID),
			},
		},
		TableName: aws.String(d.tableName),
	}

	result, err := d.table.GetItem(input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		log.Printf("No record found with session id: %s", d.sessionID)
		return []rootSchema.BaseMessageInterface{}, nil
	}

	messages := []rootSchema.BaseMessageInterface{}
	err = dynamodbattribute.UnmarshalList(result.Item["History"].L, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (d *DynamoDBChatMessageHistory) AddUserMessage(message string) error {
	msg := rootSchema.NewHumanMessage(message)
	return d.Append(msg)
}

func (d *DynamoDBChatMessageHistory) AddAIMessage(message string) error {
	msg := rootSchema.NewAIMessage(message)
	return d.Append(msg)
}

func (d *DynamoDBChatMessageHistory) Append(message rootSchema.BaseMessageInterface) error {
	messages, err := d.Messages()
	if err != nil {
		return err
	}

	messages = append(messages, message)

	av, err := dynamodbattribute.MarshalList(messages)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"SessionId": {
				S: aws.String(d.sessionID),
			},
			"History": {
				L: av,
			},
		},
		TableName: aws.String(d.tableName),
	}

	_, err = d.table.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (d *DynamoDBChatMessageHistory) Clear() error {
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"SessionId": {
				S: aws.String(d.sessionID),
			},
		},
		TableName: aws.String(d.tableName),
	}

	_, err := d.table.DeleteItem(input)
	if err != nil {
		return err
	}

	return nil
}
