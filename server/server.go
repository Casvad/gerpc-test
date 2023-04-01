package server

import (
	"context"
	"grpc-student/models"
	"grpc-student/repositories"
	"grpc-student/studentpb"
)

type Server struct {
	studentpb.UnimplementedStudentServiceServer
	repository repositories.Repository
}

func NewStudentServer(repo repositories.Repository) *Server {

	return &Server{repository: repo}
}

func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {

	student, err := s.repository.GetStudent(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &studentpb.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {

	student := &models.Student{
		Id:   req.Id,
		Name: req.Name,
		Age:  req.Age,
	}

	err := s.repository.SetStudent(ctx, student)

	if err != nil {
		return nil, err
	}

	return &studentpb.SetStudentResponse{Id: student.Id}, nil
}
