package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/specter25/go-grpc-mongo-crud/conn"
	blogpb "github.com/specter25/go-grpc-mongo-crud/protos/blog"

	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
	"github.com/specter25/go-grpc-mongo-crud/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Connect to the mongodb Instance
	conn.ConnectDB()
	db := conn.GetMongoClient()
	blogdb := db.Database("my_databse").Collection("blogs")
	//establish the grpc instance
	l := hclog.Default()
	opts := []grpc.ServerOption{}
	gs := grpc.NewServer(opts...)
	srv := server.NewBlogServiceServer(l, db, blogdb)
	blogpb.RegisterBlogServiceServer(gs, srv)
	reflection.Register(gs)

	//Start out listner at port no 50051 , this is the defaulat groc port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		l.Error("unable to listen to port no 50051")
	}

	//satrt the server
	go func() {
		err := gs.Serve(listener)
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	//handle graceful server shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Error("Received terminate , graceful shutdown", sig)
	gs.Stop()
	listener.Close()

}
