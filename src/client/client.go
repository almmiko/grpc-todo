package client

import (
	"context"
	"grpc-todo/src/pb"
	"grpc-todo/src/service"
	"log"

	"fmt"

	"google.golang.org/grpc"
)

func Boot() {
	conn, err := grpc.Dial(service.PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}

	client := pb.NewTodoServiceClient(conn)
	ctx := context.Background()

	res, err := client.DeleteTodo(ctx, &pb.TodoRequest{
		Id: "1",
	})
	if err != nil {
		log.Fatalf("get todo error: %v", err)
	}

	// res, err := client.CreateTodo(ctx, &pb.CreateTodoRequest{Completed: false, Message: "some", Title: "new todo"})
	// if err != nil {
	// 	log.Fatalf("get todo error: %v", err)
	// }

	// res, err := client.GetTodo(ctx, &pb.TodoRequest{Id: "id"})
	// if err != nil {
	// 	log.Fatalf("get todo error: %v", err)
	// }

	fmt.Println(res)
}
