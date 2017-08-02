package dynamodb

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	dockertest "gopkg.in/ory-am/dockertest.v3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func prepareDynamoDBContainer(t *testing.T) (func(), string) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("couldn't not connect docker host: %s", err.Error())
	}

	resource, err := pool.Run("atlassianlabs/localstack", "latest", []string{})
	if err != nil {
		t.Fatalf("couldn't start S3 container: %s", err.Error())
	}

	addr := fmt.Sprintf("http://localhost:%s", resource.GetPort("4569/tcp"))

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			t.Fatalf("couldn't cleanup DynamoDB container: %s", err.Error())
		}
	}

	if err = pool.Retry(func() error {
		resp, err := http.Get(addr)
		if err != nil {
			return err
		}

		if resp.StatusCode != 400 {
			return fmt.Errorf("didn't return status code 400: %s", resp.Status)
		}

		return nil
	}); err != nil {
		t.Fatalf("couldn't prepare DynamoDB container: %s", err.Error())
	}

	return cleanup, addr
}

func TestCreateTable(t *testing.T) {
	cleanup, addr := prepareDynamoDBContainer(t)
	defer cleanup()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
		Region:      aws.String(endpoints.ApNortheast1RegionID),
		Endpoint:    aws.String(addr),
	}))

	d := &DynamoDB{Service: dynamodb.New(sess)}
	if err := d.CreateTable("test"); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	svc := dynamodb.New(sess)

	expectedAttr := []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("user"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("database"),
			AttributeType: aws.String("S"),
		},
	}

	expectedSchema := []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("user"),
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("database"),
			KeyType:       aws.String("RANGE"),
		},
	}

	actual, err := svc.DescribeTable(&dynamodb.DescribeTableInput{TableName: aws.String("test")})
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	actualAttr := actual.Table.AttributeDefinitions
	actualSchema := actual.Table.KeySchema

	if !reflect.DeepEqual(expectedAttr, actualAttr) {
		t.Fatalf("didn't match AttributeDefinitions: expected %s, actual %s", expectedAttr, actualAttr)
	}

	if !reflect.DeepEqual(expectedSchema, actualSchema) {
		t.Fatalf("didn't match KeySchema: expected %s, actual %s", expectedSchema, actualSchema)
	}
}

func TestInsert(t *testing.T) {
	cleanup, addr := prepareDynamoDBContainer(t)
	defer cleanup()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
		Region:      aws.String(endpoints.ApNortheast1RegionID),
		Endpoint:    aws.String(addr),
	}))

	d := &DynamoDB{Service: dynamodb.New(sess)}
	if err := d.CreateTable("test"); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	input := map[string]string{
		"alice_db": "applepie",
		"bob_db":   "bananatart",
	}

	if err := d.Insert("test", "alice", input); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	svc := dynamodb.New(sess)

	expected := []map[string]*dynamodb.AttributeValue{
		{
			"user": {
				S: aws.String("alice"),
			},
			"database": {
				S: aws.String("alice_db"),
			},
			"password": {
				S: aws.String("applepie"),
			},
		},
		{
			"user": {
				S: aws.String("alice"),
			},
			"database": {
				S: aws.String("bob_db"),
			},
			"password": {
				S: aws.String("bananatart"),
			},
		},
	}

	param := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String("alice"),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#user": aws.String("user"),
		},
		KeyConditionExpression: aws.String("#user = :user"),
		TableName:              aws.String("test"),
	}

	actual, err := svc.Query(param)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual.Items) {
		t.Fatalf("didn't match Items: expected %s, actual %s", expected, actual.Items)
	}
}

func TestDelete(t *testing.T) {
	cleanup, addr := prepareDynamoDBContainer(t)
	defer cleanup()

	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", ""),
		Region:      aws.String(endpoints.ApNortheast1RegionID),
		Endpoint:    aws.String(addr),
	}))

	d := &DynamoDB{Service: dynamodb.New(sess)}
	if err := d.CreateTable("test"); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	input := map[string]string{
		"alice_db": "applepie",
		"bob_db":   "bananatart",
	}

	if err := d.Insert("test", "alice", input); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if err := d.Delete("test", "alice", []string{"alice_db", "bob_db"}); err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	svc := dynamodb.New(sess)

	param := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":user": {
				S: aws.String("alice"),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#user": aws.String("user"),
		},
		KeyConditionExpression: aws.String("#user = :user"),
		TableName:              aws.String("test"),
	}

	expected := []map[string]*dynamodb.AttributeValue{}

	actual, err := svc.Query(param)
	if err != nil {
		t.Fatalf("got an err: %s", err.Error())
	}

	if !reflect.DeepEqual(expected, actual.Items) {
		t.Fatalf("didn't match Items: expected %s, actual %s", expected, actual.Items)
	}
}
