package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"grpc-golang/greet/greetpb"

	"google.golang.org/grpc"
)

type server struct {}

func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {

  fmt.Println("Greet function was invoked:", req)
  firstName := req.GetGreeting().GetFirstName()
  result := "Hello "+ firstName
  res := &greetpb.GreetingResponse{
    Result: result,
  }

  return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error  {
  fmt.Printf("Greet Many Times function was invoked with %v\n\n", req)
	
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello "+ firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

	  _ = stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
  fmt.Printf("Greet function was invoked with as stream request")
  result := ""

  for  {
		req, err := stream.Recv()
		if err == io.EOF {
			// finished reading stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
  }
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with as stream request")
	
	for  {
		req, err := stream.Recv()
		if err == io.EOF {
			// finished reading stream
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "
		
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
  }
}

func main() {
	fmt.Println("Hello iam server")
	lis, err := net.Listen("tcp", "localhost:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to load server: %v", err)
	}
}