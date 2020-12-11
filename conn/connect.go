package conn

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client
var mongoCtx context.Context

func ConnectDB() {

	var err error
	log := hclog.Default()
	log.Debug("Connecting to MongoDB....")
	mongoCtx := context.Background()

	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Error("connot connect with the database", err)
	}

	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Error("Could not connect to MongoDB: %v\n", err)
	} else {
		log.Info("Connected to Mongodb")
	}
	// databases, err := db.ListDatabaseNames(mongoCtx, bson.M{})
	// if err != nil {
	// 	log.Error("error", err)
	// }
	// log.Debug("adatabases", databases)

}

func GetMongoClient() *mongo.Client {
	return db
}
