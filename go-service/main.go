package main

import (
	"context"
	"encoding/json"
	"log"
	"my-property/go-service/database"
	"my-property/go-service/graphql"
	"my-property/go-service/handlers"
	"my-property/go-service/services"
	"my-property/go-service/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"nhooyr.io/websocket"
)

func graphqlHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		result := graphql.ExecuteQuery(params.Query, params.Variables)
		c.JSON(http.StatusOK, result)
	}
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "internal error")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// 1. Read subscription query from client
	_, msg, err := c.Read(ctx)
	if err != nil {
		log.Println("Error reading subscription query:", err)
		return
	}

	var req struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}
	if err := json.Unmarshal(msg, &req); err != nil {
		log.Println("Error unmarshalling subscription query:", err)
		return
	}

	// 2. Execute the subscription
	ch := graphql.ExecuteSubscription(req.Query, req.Variables)

	// 3. Forward events to client
	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			payload, err := json.Marshal(map[string]interface{}{"data": data})
			if err != nil {
				log.Println("Error marshalling data:", err)
				continue
			}
			if err := c.Write(ctx, websocket.MessageText, payload); err != nil {
				log.Println("Error writing to websocket:", err)
				return
			}
		}
	}
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Database
	database.InitDB()

	// Initialize services
	encryptionKey := os.Getenv("ENCRYPTION_SECRET_KEY")
	if encryptionKey == "" {
		log.Fatal("ENCRYPTION_SECRET_KEY environment variable not set")
	}
	encryptionService := utils.NewEncryptionService(encryptionKey)
	financialService := services.NewFinancialService(database.DB, encryptionService)

	// Initialize GraphQL resolvers with services
	graphql.InitializeResolvers(financialService)

	router := gin.Default()

	// GraphQL endpoint
	router.POST("/graphql", graphqlHandler())
	router.GET("/graphql", graphqlHandler()) // for GraphiQL

	// REST endpoint for file upload
	router.POST("/properties/:id/upload", handlers.UploadPropertyImages)

	// Static file serving for property images
	router.Static("/storage", "./storage")

	// Use gin.WrapF to convert the http.HandlerFunc to a gin.HandlerFunc
	router.GET("/ws", func(c *gin.Context) {
		websocketHandler(c.Writer, c.Request)
	})

	// Serve the application
	log.Println("Go server is running on port 8080")
	log.Fatal(router.Run(":8080"))
}
