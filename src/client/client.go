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

	res, err := client.GetTodo(ctx, &pb.TodoRequest{Id: "id"})
	if err != nil {
		log.Fatalf("get todo error: %v", err)
	}

	fmt.Println(res)
}
