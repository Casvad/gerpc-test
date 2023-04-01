package server

import (
	"context"
	"grpc-student/models"
	"grpc-student/repositories"
	"grpc-student/studentpb"
	"grpc-student/testpb"
	"io"
	"log"
	"time"
)

type TestServer struct {
	testpb.UnimplementedTestServiceServer
	repository repositories.Repository
}

func NewTestServer(repo repositories.Repository) *TestServer {

	return &TestServer{repository: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {

	test, err := s.repository.GetTest(ctx, req.Id)

	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {

	test := &models.Test{
		Id:   req.Id,
		Name: req.Name,
	}

	err := s.repository.SetTest(ctx, test)

	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{Id: test.Id}, nil
}

func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}

		question := models.Question{
			Id:       msg.Id,
			Question: msg.Question,
			Answer:   msg.Answer,
			TestId:   msg.TestId,
		}

		err = s.repository.SetQuestion(context.Background(), &question)

		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}

		enrollment := models.Enrollment{
			StudentId: msg.StudentId,
			TestId:    msg.TestId,
		}

		err = s.repository.SetEnrollment(context.Background(), &enrollment)
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(request *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {

	students, err := s.repository.GetStudentsPerTest(context.Background(), request.TestId)

	if err != nil {

		return err
	}

	for _, student := range students {
		studentPb := studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}

		err = stream.Send(&studentPb)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {

	questions, err := s.repository.GetQuestionsPerTest(context.Background(), "t1")
	if err != nil {
		return err
	}
	i := 0
	var currentQuestion = &models.Question{}
	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}
		if i <= len(questions) {
			questionToSend := &testpb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err = stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
		}
		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return nil
		}
		log.Println("Answer: ", answer.GetAnswer())
	}
}
