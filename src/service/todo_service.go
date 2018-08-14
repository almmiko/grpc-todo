package service

import (
	"context"
	"grpc-todo/src/pb"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TodoService struct{}

var dbConn *dynamodb.DynamoDB

const (
	TABLE = "Todos"
	PORT  = ":9090"
)

func (t *TodoService) GetTodo(ctx context.Context, req *pb.TodoRequest) (*pb.TodoResponse, error) {

	todo := pb.Todo{}

	filter := expression.Name("id").Equal(expression.Value(req.Id))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		log.Fatalf("Got error building expression: %v", err.Error())
		return nil, err
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(TABLE),
	}

	result, err := dbConn.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %v", err.Error())
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, status.Errorf(codes.NotFound,
			"Todo with id %s not found", req.GetId())
	}

	for _, item := range result.Items {
		if err := dynamodbattribute.UnmarshalMap(item, &todo); err != nil {
			log.Fatalf("UnmarshalMap failed: %v", err.Error())
			return nil, err
		}
	}

	return &pb.TodoResponse{
		Todo: &todo,
	}, nil
}

func (t *TodoService) GetTodos(ctx context.Context, req *pb.TodosRequest) (*pb.TodosResponse, error) {

	var todos []*pb.Todo

	params := &dynamodb.ScanInput{
		TableName: aws.String(TABLE),
	}

	result, err := dbConn.Scan(params)
	if err != nil {
		log.Fatalf("Query API call failed: %v", err.Error())
		return nil, err
	}

	for _, item := range result.Items {
		todo := pb.Todo{}

		if err := dynamodbattribute.UnmarshalMap(item, &todo); err != nil {
			log.Fatalf("UnmarshalMap failed: %v", err.Error())
			return nil, err
		}

		todos = append(todos, &todo)
	}

	return &pb.TodosResponse{
		Todos: todos,
	}, nil
}

func (t *TodoService) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.TodoResponse, error) {

	uuId := uuid.NewV4()

	todo := &pb.Todo{
		Id:        uuId.String(),
		Completed: req.Completed,
		Title:     req.Title,
		Message:   req.Message,
	}

	av, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		log.Fatalf("Error marshalling map: %v", err.Error())
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(TABLE),
	}

	_, err = dbConn.PutItem(input)
	if err != nil {
		log.Fatalf("Error calling PutItem: %v", err.Error())
		return nil, err
	}

	return &pb.TodoResponse{
		Todo: todo,
	}, nil
}

func (t *TodoService) UpdateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoResponse, error) {

	updatedTodo := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {
				S: aws.String(req.Title),
			},
			":m": {
				S: aws.String(req.Message),
			},
			":c": {
				BOOL: aws.Bool(req.Completed),
			},
		},
		TableName:        aws.String(TABLE),
		UpdateExpression: aws.String("set title = :t, message = :m, completed = :c"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(req.Id),
			},
		},
		ReturnValues: aws.String("ALL_NEW"),
	}

	result, err := dbConn.UpdateItem(updatedTodo)
	if err != nil {
		log.Fatalf("Error calling UpdateItem: %v", err.Error())
		return nil, err
	}

	todo := pb.Todo{}

	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &todo); err != nil {
		log.Fatalf("Error calling UnmarshalMap: %v", err.Error())
		return nil, err
	}

	return &pb.TodoResponse{
		Todo: &todo,
	}, nil
}

func (t *TodoService) DeleteTodo(ctx context.Context, req *pb.TodoRequest) (*pb.TodoResponse, error) {
	deletedTodo := &dynamodb.DeleteItemInput{
		TableName: aws.String(TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(req.Id),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
	}

	result, err := dbConn.DeleteItem(deletedTodo)
	if err != nil {
		log.Fatalf("Error calling DeleteItem: %v", err.Error())
		return nil, err
	}

	todo := pb.Todo{}

	if err := dynamodbattribute.UnmarshalMap(result.Attributes, &todo); err != nil {
		log.Fatalf("Error calling UnmarshalMap: %v", err.Error())
		return nil, err
	}

	return &pb.TodoResponse{
		Todo: &todo,
	}, nil
}

func Boot() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	log.Println("TodoService is listening on ", PORT)

	server := grpc.NewServer()

	pb.RegisterTodoServiceServer(server, &TodoService{})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

func NewDatabaseConnection() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)

	return dynamodb.New(sess), err
}

func init() {
	conn, err := NewDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed database connection %v", err)
	}

	dbConn = conn
}
