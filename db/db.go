package db

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Johnhi19/TreeSpotter_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func Connect(txtFile string) {
	file, err := os.Open(txtFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.rpy40.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", username, password)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func InsertOneMeadow(meadow models.Meadow) interface{} {
	Connect("credentials.txt")

	collection := client.Database("TreeSpotter").Collection("Meadow")

	// Convert Meadow struct to BSON document
	meadowDoc := bson.D{
		{"ID", meadow.ID},
		{"Name", meadow.Name},
		{"NumberTrees", meadow.NumberTrees},
	}

	result, err := collection.InsertOne(context.TODO(), meadowDoc)

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted a meadow with ID:", result.InsertedID)
	Disconnect()

	return result.InsertedID
}

func InsertOneTree(tree bson.D) {
	Connect("credentials.txt")

	collection := client.Database("TreeSpotter").Collection("Tree")

	result, err := collection.InsertOne(context.TODO(), tree)

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted a meadow with ID:", result.InsertedID)
	Disconnect()

}

func FindOneMeadowById(filter bson.D) bson.M {
	Connect("credentials.txt")

	collection := client.Database("TreeSpotter").Collection("Meadow")
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found meadow:", result)
	Disconnect()

	return result
}

func FindOneTree(filter bson.D) {
	Connect("credentials.txt")

	collection := client.Database("TreeSpotter").Collection("Tree")
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found tree:", result)
	Disconnect()
}

func Disconnect() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	fmt.Println("Disconnected from MongoDB!")
}
