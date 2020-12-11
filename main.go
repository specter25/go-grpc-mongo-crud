package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/specter25/go-grpc-mongo-crud/conn"
	blogpb "github.com/specter25/go-grpc-mongo-crud/protos/blogs"

	"github.com/hashicorp/go-hclog"
	"github.com/specter25/go-grpc-mongo-crud/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	l := hclog.Default()

	//Set options , here we can configure things like TLS support
	opts := []grpc.ServerOption{}
	//Create a new grpc server with blank support
	gs := grpc.NewServer(opts...)
	//create blogs type
	srv := server.NewBlogs(l)
	//Register the service with the server
	blogpb.RegisterBlogServiceServer(gs, srv)
	reflection.Register(gs)

	//Start out listner at port no 50051 , this is the defaulat groc port
	listener, err := net.Listen("tcp", ":50051")
	//Handle any errors if they occur
	if err != nil {
		l.Error("unable to listena t port no 50051")
	}

	//Connect to the mongodb Instance
	conn.ConnectDB()

	go func() {
		err := gs.Serve(listener)
		if err != nil {
			l.Error("Error starting server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)
	sig := <-sigChan
	l.Error("Received terminate , graceful shutdown", sig)
	gs.Stop()
	listener.Close()

}
