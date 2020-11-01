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
	// doClientStreaming(c)
	doBiDiStreaming(c)
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
		sendErr := stream.Send(req)
		
		if sendErr != nil {
			log.Fatalf("Error sending data: %v", err)
		}
		
	  time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while recieving response from long greet: %v", err)
	}
	fmt.Printf("LongGreet Response : %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient)  {
  fmt.Println("Starting to do a BiDi Streaming RPC..")
	
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling BiDi: %v", err)
		return
	}
	
	requests := []*greetpb.GreetEveryoneRequest{
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
	
	waitc := make(chan struct{})
	// send bunch message to client (go routine)
	go func(){
		//function to send bunch message to client
		for _, req := range requests {
			fmt.Printf("Sending Message: %v\n", req)
			sendErr := stream.Send(req)
		
			if sendErr != nil {
				log.Fatalf("Error sending data: %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		
		stream.CloseSend()
	}()
	
	// recieve bunch message to client (go routine)
	go func(){
		//function to recieve bunch message from client
		for {
			res, err := stream.Recv()
			
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while recieve: %v", err)
				break
			}
			
			fmt.Printf("Recieved: %v\n", res.GetResult())
		}
		close(waitc)
		
	}()
	
	//block until everything is done
	<-waitc
}
