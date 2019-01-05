package database

import (
	"context"
	"github.com/curiouscat2018/helloworld-api/common"

	"github.com/gin-gonic/gin"
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

func (db azureDatabase) GetEntry(c *gin.Context) (*Entry, error) {
	filter := bson.M{"greeting": "helloworld!"}
	update := bson.M{"$inc": bson.M{"request_count": 1}}

	result := db.collection.FindOneAndUpdate(c, filter, update)
	entry := Entry{}

	if err := result.Decode(&entry); err != nil {
		common.TraceWarn(c).Err(err).Msg("db record not found or corrupted")

		entry.Greeting = "helloworld!"
		entry.RequestCount = 1
		_, err := db.collection.InsertOne(c, entry)

		if err != nil {
			return nil, err
		}
	}
	return &entry, nil
}
