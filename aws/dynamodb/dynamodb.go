package dynamodb

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	awspkg "github.com/munisystem/invoker/aws"
)

type DynamoDB struct {
	Service *dynamodb.DynamoDB
}

func NewClient() *DynamoDB {
	return &DynamoDB{
		Service: dynamodb.New(awspkg.Session()),
	}
}

func (d *DynamoDB) CreateTable(name string) error {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("user"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("database"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("user"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("database"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(name),
	}

	if _, err := d.Service.CreateTable(input); err != nil {
		return fmt.Errorf("Faild to create DynamoDB table: %s", err.Error())
	}

	return nil
}

func (d *DynamoDB) Insert(table, user string, items map[string]string) error {
	var writeRequests []*dynamodb.WriteRequest

	for key, value := range items {
		wr := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: map[string]*dynamodb.AttributeValue{
					"user": {
						S: aws.String(user),
					},
					"database": {
						S: aws.String(key),
					},
					"password": {
						S: aws.String(value),
					},
				},
			},
		}

		writeRequests = append(writeRequests, wr)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{table: writeRequests},
	}

	if _, err := d.Service.BatchWriteItem(input); err != nil {
		return fmt.Errorf("Faild to insert items into DynamoDB (table: %s): %s", table, err.Error())
	}

	return nil
}

func (d *DynamoDB) Delete(table, user string, databases []string) error {
	var deleteRequests []*dynamodb.WriteRequest

	for _, database := range databases {
		dr := &dynamodb.WriteRequest{
			DeleteRequest: &dynamodb.DeleteRequest{
				Key: map[string]*dynamodb.AttributeValue{
					"user": {
						S: aws.String(user),
					},
					"database": {
						S: aws.String(database),
					},
				},
			},
		}

		deleteRequests = append(deleteRequests, dr)
	}

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{table: deleteRequests},
	}

	if _, err := d.Service.BatchWriteItem(input); err != nil {
		return fmt.Errorf("Faild to delete items from DynamoDB (table: %s): %s", table, err.Error())
	}

	return nil
}

func (d *DynamoDB) Get(table, user string) (map[string]string, error) {
	param := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String(user),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#user": aws.String("user"),
		},
		KeyConditionExpression: aws.String("#user = :user"),
		TableName:              aws.String("test"),
	}

	result, err := d.Service.Query(param)
	if err != nil {
		return nil, fmt.Errorf("Faild to get items from DynamoDB (table: %s): %s", table, err.Error())
	}

	items := make(map[string]string)

	for _, item := range result.Items {
		database := *item["database"].S
		password := *item["password"].S
		items[database] = password
	}

	return items, nil
}
