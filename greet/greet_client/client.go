package main

import (
  "context"
  "fmt"
	"google.golang.org/grpc"
	"grpc-golang/greet/greetpb"
  "io"
  "log"
)

func main() {
	fmt.Println("Hello client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connected: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	fmt.Printf("Create client : %f", c)
	req := &greetpb.GreetingRequest{
	  Greeting: &greetpb.Greeting{
	    FirstName: "Adam Lesmana",
	    LastName: "Ganda Saputra",
	  },
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil{
	  log.Fatalf("Error calling Greet Rpc: %v", err)
	}

	log.Printf("Response from guets: %v", res.Result)

	//doUnary(c)

	// Server Streaming
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient)  {
  fmt.Println("Starting to do a Unary RPC..")
  req := &greetpb.GreetingRequest{
	Greeting: &greetpb.Greeting{
	  FirstName: "Gavin Khawarimi",
	  LastName: "Abqary",
	},
  }
  res, err := c.Greet(context.Background(), req)

  if err != nil{
	log.Fatalf("Error calling Greet Rpc: %v", err)
  }

  log.Printf("Response from guets: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient)  {
  fmt.Println("Starting to do a Server Streaming RPC..")
  req := &greetpb.GreetManyTimesRequest{
	Greeting: &greetpb.Greeting{
	  FirstName: "Adam Lesmana",
	  LastName: "Ganda Saputra",
	},
  }
  resStream, err := c.GreetManyTimes(context.Background(), req)

  if err != nil{
	log.Fatalf("Error calling Server Streaming Greet Many Times Rpc: %v", err)
  }

  for  {
	msg, err := resStream.Recv()
	if err == io.EOF {
		// end of the stream
	  break;
	}

	if err !=  nil {
	  log.Fatalf("Error reading stream: %v", err)
	}

	log.Printf("Response from GreatManyTimes: %v", msg.GetResult())
  }
}