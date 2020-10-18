package main

import (
  "context"
  "fmt"
	"log"
	"net"
	"strconv"
  "time"

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

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error  {
  fmt.Println("Greet Many Times function was invoked with %v\n", req)
	
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello "+ firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	
	return nil
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