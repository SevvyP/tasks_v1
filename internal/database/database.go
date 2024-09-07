package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Task struct {
	ID    string `json:"id" dynamodbav:"id"`
	Items []Item `json:"items" dynamodbav:"items"`
}

type Item struct {
	Name  string `json:"name" dynamodbav:"name"`
	Price int    `json:"price" dynamodbav:"price"`
}

// TaskDatabase is an interface for interacting with a task database.
type TaskDatabase interface {
	GetTasks() (*[]Task, error)
	GetTaskByID(id string) (*Task, error)
	CreateTask(task Task) error
	UpdateTask(task Task) error
	DeleteTask(task Task) error
}

// DynamoDatabase is a DynamoDB implementation of the TaskDatabase interface.
type DynamoDatabase struct {
	tableName string
	client    *dynamodb.Client
}

// NewDatabase creates a new DynamoDatabase instance.
// It returns an error if the AWS configuration cannot be loaded.
func NewDatabase() (*DynamoDatabase, error) {
	tableName := "tasks_v1"
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %v", err)
	}

	return &DynamoDatabase{
		tableName: tableName,
		client:    dynamodb.NewFromConfig(cfg),
	}, nil
}

// GetTasks retrieves a list of tasks from the database.
// It returns an error if the tasks cannot be retrieved.
func (d *DynamoDatabase) GetTasks() (*[]Task, error) {
	// Create an input for the Scan operation
	input := &dynamodb.ScanInput{
		TableName: aws.String(d.tableName),
	}

	// Perform the Scan operation
	out, err := d.client.Scan(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}

	output := make([]Task, len(out.Items))
	err = attributevalue.UnmarshalListOfMaps(out.Items, &output)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}

	// Replace resultresult.Items with result.Items
	return &output, nil
}

// GetTaskByID retrieves a task by its ID from the database.
// It returns an error if the task cannot be retrieved.
// If the task is not found, it returns nil.
func (d *DynamoDatabase) GetTaskByID(id string) (*Task, error) {
	// Marshal the id into a DynamoDB attribute value
	idAttr, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal id: %v", err)
	}

	// Create an input for the GetItem operation
	input := &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id": idAttr,
		},
	}

	// Perform the GetItem operation
	out, err := d.client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	var task Task
	err = attributevalue.UnmarshalMap(out.Item, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	if task.ID == "" {
		return nil, nil
	}

	return &task, nil
}

// CreateTask creates a new task in the database.
// It returns an error if the task cannot be created.
func (d *DynamoDatabase) CreateTask(task Task) error {
	// Marshal the task into a DynamoDB attribute value map
	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	// Create an input for the PutItem operation
	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	}

	// Perform the PutItem operation
	_, err = d.client.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

// UpdateTask updates an existing task in the database.
// It returns an error if the task cannot be updated.
func (d *DynamoDatabase) UpdateTask(updatedTask Task) error {
	id, err := attributevalue.Marshal(updatedTask.ID)
	if err != nil {
		return fmt.Errorf("failed to marshal updated task: %v", err)
	}

	// Check if the task exists
	exists, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id": id,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to check for task: %v", err)
	}
	if len(exists.Item) == 0 {
		return fmt.Errorf("task does not exist")
	}

	// Marshal the task into a DynamoDB attribute value map
	item, err := attributevalue.MarshalMap(updatedTask)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %v", err)
	}

	// Create an input for the PutItem operation
	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	}

	// Perform the PutItem operation
	_, err = d.client.PutItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

// DeleteTask deletes a task from the database.
// It returns an error if the task cannot be deleted.
func (d *DynamoDatabase) DeleteTask(taskToDelete Task) error {
	id, err := attributevalue.Marshal(taskToDelete.ID)
	if err != nil {
		return fmt.Errorf("failed to marshal updated task: %v", err)
	}

	// Create an input for the DeleteItem operation
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"id": id,
		},
	}

	// Perform the DeleteItem operation
	_, err = d.client.DeleteItem(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}
