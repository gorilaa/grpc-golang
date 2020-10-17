package main

import (
  "context"
  "fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"grpc-golang/greet/greetpb"
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