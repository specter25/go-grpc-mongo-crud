package main

import (
	"net"
	"os"
	"os/signal"

	"github.com/specter25/go-grpc-mongo-crud/conn"
	blogpb "github.com/specter25/go-grpc-mongo-crud/protos/blog"

	"github.com/hashicorp/go-hclog"
	"github.com/specter25/go-grpc-mongo-crud/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	//establish the grpc instance
	l := hclog.Default()
	opts := []grpc.ServerOption{}
	gs := grpc.NewServer(opts...)
	srv := server.NewBlogServiceServer(l)
	blogpb.RegisterBlogServiceServer(gs, srv)
	reflection.Register(gs)

	//Start out listner at port no 50051 , this is the defaulat groc port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		l.Error("unable to listen to port no 50051")
	}

	//Connect to the mongodb Instance
	conn.ConnectDB()

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
