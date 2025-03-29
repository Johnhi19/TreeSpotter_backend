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

func FindAllMeadows() []bson.M {
	collection := client.Database("TreeSpotter").Collection("Meadow")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	var meadows []bson.M
	if err = cursor.All(context.Background(), &meadows); err != nil {
		panic(err)
	}

	fmt.Println("Found meadows:", meadows)
	return meadows
}

func FindAllTreesForMeadow(filter bson.D) []bson.M {
	collection := client.Database("TreeSpotter").Collection("Tree")
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var trees []bson.M
	if err = cursor.All(context.Background(), &trees); err != nil {
		panic(err)
	}

	fmt.Println("Found trees:", trees)
	return trees
}

func InsertOneMeadow(meadow models.Meadow) interface{} {
	collection := client.Database("TreeSpotter").Collection("Meadow")

	meadowDoc := models.TransformMeadowToBson(meadow)

	result, err := collection.InsertOne(context.TODO(), meadowDoc)

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted a meadow with ID:", result.InsertedID)
	return result.InsertedID
}

func InsertOneTree(tree models.Tree) interface{} {
	collection := client.Database("TreeSpotter").Collection("Tree")

	treeDoc := models.TransformTreeToBson(tree)

	result, err := collection.InsertOne(context.TODO(), treeDoc)

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted a meadow with ID:", result.InsertedID)
	return result.InsertedID
}

func FindOneMeadowById(filter bson.D) bson.M {
	collection := client.Database("TreeSpotter").Collection("Meadow")
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found meadow:", result)
	return result
}

func FindOneTree(filter bson.D) bson.M {
	collection := client.Database("TreeSpotter").Collection("Tree")
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found tree:", result)
	return result
}

func Disconnect() {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
	fmt.Println("Disconnected from MongoDB!")
}
