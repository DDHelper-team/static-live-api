package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	//將live.live裡面的東西儲存到output.json

	live := client.Database("live").Collection("live")
	cursor, err := live.Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	file, err := os.Create("live.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	list_result := []bson.M{}
	for _, result := range results {
		delete(result, "_id")
		list_result = append(list_result, result)
	}
	encoder.Encode(list_result)
}
