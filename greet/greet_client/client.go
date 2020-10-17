package main

import (
	"fmt"
	"google.golang.org/grpc"
	"grpc-golang/greet/greetpb"
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
}