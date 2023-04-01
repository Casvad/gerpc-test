package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-student/testpb"
	"io"
	"log"
	"time"
)

func main() {

	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("could not connnect: %v", err)
	}
	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)

	//DoUnary(c)
	//DoClientStreaming(c)
	//DoServerStreaming(c)
	DoBidirectionalStreaming(c)
}

func DoUnary(c testpb.TestServiceClient) {

	req := &testpb.GetTestRequest{
		Id: "t1",
	}

	res, err := c.GetTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Error calling getTest: %v", err)
	}

	log.Printf("response from server %v", res)
}

func DoClientStreaming(c testpb.TestServiceClient) {

	questions := []*testpb.Question{
		{
			Id:       "q8",
			Answer:   "t1",
			Question: "azul",
			TestId:   "t1",
		},
		{
			Id:       "q9",
			Answer:   "t1",
			Question: "azul",
			TestId:   "t1",
		},
		{
			Id:       "q10",
			Answer:   "Especialidad de golang",
			Question: "backend",
			TestId:   "t1",
		},
	}

	stream, err := c.SetQuestions(context.Background())
	if err != nil {
		log.Fatalf("Error while calling setQuestions: %v", err)
	}
	for _, question := range questions {
		log.Printf("Sending question: %s \n", question.Id)
		err = stream.Send(question)
		if err != nil {
			log.Fatalf("Error sending data %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receiving response %v", err)
	}

	log.Printf("Success sending mesages %v", msg)
}

func DoServerStreaming(c testpb.TestServiceClient) {

	req := &testpb.GetStudentsPerTestRequest{TestId: "t1"}

	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Error 1 %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error 2 %v", err)
		}
		log.Printf("Received: %v", msg)
	}
}

func DoBidirectionalStreaming(c testpb.TestServiceClient) {

	answer := testpb.TakeTestRequest{
		Answer: "42",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())

	if err != nil {
		log.Fatalf("Error 1 %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			err := stream.Send(&answer)
			if err != nil {
				log.Fatalf("Error 2 %v", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error 3 %v", err)
				break
			}
			log.Printf("Response: %v", res)
		}
		close(waitChannel)
	}()
	<-waitChannel
}
