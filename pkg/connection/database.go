package connection

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func ConnectToMongo(URI string) (*mongo.Client, error) {
	fmt.Println("trying to connect to mongo db")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().
		ApplyURI(URI))

	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		var ConnectionError = fmt.Errorf("Error connecting to MongoDB: %w", err)
		return nil, ConnectionError
	}
	// fmt.Println(fmt.Printf("%+v\n", client))
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}
	fmt.Println("Connected to MongoDB!")
	Client = client
	return Client, nil
}
