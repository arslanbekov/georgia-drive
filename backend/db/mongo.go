package db

import (
	"context"
	"github.com/sirupsen/logrus" // Импорт logrus
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

var client *mongo.Client

func init() {
	var err error

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logrus.Fatal("MONGO_URI environment variable is not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatal("Failed to connect to MongoDB:", err)
	} else {
		logrus.Info("Successfully connected to MongoDB")
	}
}

func SaveToMongo(collectionName string, data map[string]interface{}) {
	logrus.Infof("Saving data to MongoDB collection: %s", collectionName)
	collection := client.Database("dates").Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		logrus.Error("Failed to insert data into MongoDB:", err)
	} else {
		logrus.Infof("Data successfully saved to collection: %s", collectionName)
	}
}

func ClearCollection(collectionName string) {
	logrus.Infof("Clearing data in MongoDB collection: %s", collectionName)
	collection := client.Database("dates").Collection(collectionName)
	_, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		logrus.Error("Failed to clear MongoDB collection:", err)
	} else {
		logrus.Infof("Collection %s cleared successfully", collectionName)
	}
}

func FetchFromMongo(collectionName string) []map[string]interface{} {
	logrus.Infof("Fetching data from MongoDB collection: %s", collectionName)
	collection := client.Database("dates").Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		logrus.Error("Failed to fetch data from MongoDB:", err)
		return nil
	}
	defer cursor.Close(context.TODO())

	var results []map[string]interface{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		logrus.Error("Failed to decode results from MongoDB:", err)
		return nil
	}
	logrus.Infof("Successfully fetched data from collection: %s", collectionName)
	return results
}

func SaveExecutionTime(executionTime time.Time) {
	collection := client.Database("dates").Collection("last_exec_time")
	_, err := collection.InsertOne(context.TODO(), map[string]interface{}{
		"timestamp": executionTime,
	})
	if err != nil {
		logrus.Error("Failed to save execution time:", err)
	}
}
