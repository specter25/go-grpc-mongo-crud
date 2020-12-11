package conn

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client
var mongoCtx context.Context

func ConnectDB() {
	log := hclog.Default()
	log.Debug("Connecting to MongoDB....")
	mongoCtx := context.Background()
	//Connect takes in a context and options , the connection uri is the only option we pass now
	db, err := mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb+srv://ujjwal:ujjwal@cluster1.khf9x.mongodb.net/test?retryWrites=true&w=majority"))
	// Handle potential errors
	if err != nil {
		log.Error("connot connect with the database", err)
	}
	// Check whether the connection was succesful by pinging the MongoDB server
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Error("Could not connect to MongoDB: %v\n", err)
	} else {
		log.Info("Connected to Mongodb")
	}
}

func GetMongoClient() *mongo.Client {
	return db
}
