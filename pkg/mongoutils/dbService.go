package mongoutils

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type DBService struct {
	Client     *mongo.Client
	Context    context.Context
	DB         *mongo.Database
	Collection *mongo.Collection
}

func setupConnection(uri string) (*mongo.Client, context.Context) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	context := context.Background()
	return client, context
}

func (s *DBService) PopulateIndex(key string) {
	c := s.Collection
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	index := yieldIndexModel(key)
	c.Indexes().CreateOne(context.Background(), index, opts)
	log.Println("Successfully created the index")
}

func yieldIndexModel(key string) mongo.IndexModel {
	// the Value is super wonky here, but this is how you set indexing to ascending (1) or descending (-1)
	keys := bsonx.Doc{{Key: key, Value: bsonx.Int32(int32(1))}}
	index := mongo.IndexModel{}
	index.Keys = keys

	index.Options = options.Index().SetUnique(false)

	return index
}

func ListIndexes(client *mongo.Client, database, collection string) {
	c := client.Database(database).Collection(collection)
	duration := 10 * time.Second
	batchSize := int32(10)
	cur, err := c.Indexes().List(context.Background(), &options.ListIndexesOptions{&batchSize, &duration})
	if err != nil {
		log.Fatalf("Something went wrong listing %v", err)
	}
	for cur.Next(context.Background()) {
		index := bson.D{}
		cur.Decode(&index)
		log.Println(fmt.Sprintf("index found %v", index))
	}
}

func NewDBService(address string, dbname string, collname string) *DBService {
	client, context := setupConnection(address)
	db := client.Database(dbname)
	collection := db.Collection(collname)
	return &DBService{Context: context, Collection: collection, DB: db, Client: client}
}
