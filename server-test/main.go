package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-student/database"
	server "grpc-student/server"
	"grpc-student/testpb"
	"log"
	"net"
)

func main() {

	list, err := net.Listen("tcp", ":5070")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://docker:docker@localhost:5432/student?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	serv := server.NewTestServer(repo)

	s := grpc.NewServer()
	testpb.RegisterTestServiceServer(s, serv)

	reflection.Register(s)

	if err = s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
