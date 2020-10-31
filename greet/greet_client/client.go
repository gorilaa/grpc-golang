package main

import (
	"context"
	"fmt"
	"grpc-golang/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
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
	//doServerStreaming(c)

	// Client Streaming
	doClientStreaming(c)
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
		break
	}

	if err !=  nil {
		log.Fatalf("Error reading stream: %v", err)
	}

	log.Printf("Response from GreatManyTimes: %v", msg.GetResult())
  }
}

func doClientStreaming(c greetpb.GreetServiceClient)  {
	fmt.Println("Starting client streaming RPC ...")

	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Adam Lesmana",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Robby Purwa",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Eko Kurniawan",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gavin Khawarizmi",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling long greet: %v", err)
	}

	for _, req := range requests{
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
	  time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while recieving response from long greet: %v", err)
	}
	fmt.Printf("LongGreet Response : %v\n", res)
}