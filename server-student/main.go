package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-student/database"
	server "grpc-student/server"
	"grpc-student/studentpb"
	"log"
	"net"
)

func main() {

	list, err := net.Listen("tcp", ":5060")

	if err != nil {
		log.Fatal(err)
	}

	repo, err := database.NewPostgresRepository("postgres://docker:docker@localhost:5432/student?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	serv := server.NewStudentServer(repo)

	s := grpc.NewServer()
	studentpb.RegisterStudentServiceServer(s, serv)

	reflection.Register(s)

	if err = s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
