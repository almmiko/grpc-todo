package api

import (
	"context"
	"grpc-todo/src/pb"
	"grpc-todo/src/service"
	"log"

	"net/http"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

const (
	PORT = ":8081"
)

func Boot() {
	r := gin.Default()

	conn, err := grpc.Dial(service.PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial error: %v", err)
	}

	client := pb.NewTodoServiceClient(conn)
	ctx := context.Background()

	r.GET("/todo", func(c *gin.Context) {

		res, err := client.GetTodos(ctx, &pb.TodosRequest{})
		if err != nil {
			log.Fatalf("gRPC get todos error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	})

	r.GET("/todo/:id", func(c *gin.Context) {
		todoId := c.Param("id")

		res, err := client.GetTodo(ctx, &pb.TodoRequest{
			Id: todoId,
		})
		if err != nil {
			log.Fatalf("gRPC get todo error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	})

	r.POST("/todo", func(c *gin.Context) {

		var todo pb.CreateTodoRequest

		decoder := json.NewDecoder(c.Request.Body)

		if err := decoder.Decode(&todo); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"errors": err.Error(),
			})

			return
		}

		res, err := client.CreateTodo(ctx, &todo)
		if err != nil {
			log.Fatalf("gRPC create todo error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": err.Error(),
			})

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	})

	r.PUT("/todo", func(c *gin.Context) {
		var todo pb.Todo

		decoder := json.NewDecoder(c.Request.Body)

		if err := decoder.Decode(&todo); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"errors": err.Error(),
			})

			return
		}

		res, err := client.UpdateTodo(ctx, &todo)
		if err != nil {
			log.Fatalf("gRPC update todo error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": err.Error(),
			})

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data": res,
		})

	})

	r.DELETE("/todo/:id", func(c *gin.Context) {
		todoId := c.Param("id")

		res, err := client.DeleteTodo(ctx, &pb.TodoRequest{
			Id: todoId,
		})
		if err != nil {
			log.Fatalf("gRPC delete todo error: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": res,
		})
	})

	r.Run(PORT)
}
