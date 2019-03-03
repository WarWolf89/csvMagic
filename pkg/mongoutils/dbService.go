package mongoutils

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CsvService struct {
	Context    context.Context
	Collection *mongo.Collection
}

func SetupConnection(uri string) (*mongo.Client, context.Context) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	context := context.Background()
	return client, context
}

func CreateCsvService(client *mongo.Client, db string, collName string) *CsvService {

	collection := client.Database(db).Collection(collName)
	return &CsvService{Context: context.Background(), Collection: collection}
}
