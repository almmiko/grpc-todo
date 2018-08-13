package service

import (
	"context"
	"grpc-todo/src/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type TodoService struct{}

var Todos []pb.Todo = []pb.Todo{
	pb.Todo{
		Id:        "1",
		Completed: false,
		Title:     "first todo",
		Message:   "additional message",
	},
	pb.Todo{
		Id:        "2",
		Completed: true,
		Title:     "2 todo",
		Message:   "additional message",
	},
	pb.Todo{
		Id:        "3",
		Completed: false,
		Title:     "3 todo",
		Message:   "additional message",
	},
}

func (t *TodoService) GetTodo(ctx context.Context, req *pb.TodoRequest) (*pb.TodoResponse, error) {
	todo := Todos[0]
	// if !todo {
	// 	return nil, status.Errorf(codes.NotFound,
	// 		"Todo with id %s not found", req.GetId())
	// }

	return &pb.TodoResponse{
		Todo: &todo,
	}, nil
}

func (t *TodoService) GetTodos(ctx context.Context, req *pb.TodosRequest) (*pb.TodosResponse, error) {
	return nil, nil
}

func (t *TodoService) CreateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoResponse, error) {
	return nil, nil
}

func (t *TodoService) UpdateTodo(ctx context.Context, req *pb.Todo) (*pb.TodoResponse, error) {
	return nil, nil
}

func (t *TodoService) DeleteTodo(ctx context.Context, req *pb.TodoRequest) (*pb.TodoResponse, error) {
	return nil, nil
}

const (
	PORT = ":9090"
)

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
