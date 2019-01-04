package database

import (
	"context"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type azureDatabase struct {
	collection *mongo.Collection
}

func NewAzureDatabase(url string) (Database, error) {
	db := azureDatabase{}
	client, err := mongo.Connect(context.TODO(), url)
	if err != nil {
		return nil, err
	}
	db.collection = client.Database("main").Collection("main")
	return db, nil
}

func (db azureDatabase) GetEntry() (*Entry, error) {
	filter := bson.M{"greeting": "helloworld!"}
	update := bson.M{"$inc": bson.M{"requestcount": 1}}

	result := db.collection.FindOneAndUpdate(context.TODO(), filter, update)
	entry := Entry{}
	if err := result.Decode(&entry); err != nil {
		log.Printf("db record not found or corrupted: %v\n", err)

		entry.Greeting = "helloworld!"
		entry.RequestCount = 1
		result, err := db.collection.InsertOne(context.TODO(), &entry)

		if err != nil {
			return nil, err
		}
		log.Printf("inserted Entry %v\n", result.InsertedID)
	}
	return &entry, nil
}
