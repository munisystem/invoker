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
				AttributeName: aws.String("name"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("database"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("name"),
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
