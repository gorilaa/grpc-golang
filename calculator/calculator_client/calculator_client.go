package main

import (
  "context"
  "fmt"
  "google.golang.org/grpc"
  "grpc-golang/calculator/calculatorpb"
  "log"
)

func main() {
  fmt.Println("Hello client")

  cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

  if err != nil {
	log.Fatalf("could not connected: %v", err)
  }

  defer cc.Close()

  c := calculatorpb.NewCalculatorServiceClient(cc)

  fmt.Printf("Create client : %f", c)
  req := &calculatorpb.CalculatorRequest{
	Calculator: &calculatorpb.Calculator{
	  ParameterOne: 10,
	  ParameterTwo: 22,
	},
  }
  res, err := c.Sum(context.Background(), req)

  if err != nil{
	log.Fatalf("Error calling Greet Rpc: %v", err)
  }

  log.Printf("Response from guets: %v", res.CalculatorResult)
  doSumUnary(c)
}

func doSumUnary(c calculatorpb.CalculatorServiceClient)  {
  fmt.Println("Starting to do a Unary RPC..")
  req := &calculatorpb.CalculatorRequest{
	Calculator: &calculatorpb.Calculator{
	  ParameterOne: 40,
	  ParameterTwo: 60,
	},
  }
  res, err := c.Sum(context.Background(), req)

  if err != nil{
	log.Fatalf("Error calling Greet Rpc: %v", err)
  }

  log.Printf("Response from guets: %v", res.CalculatorResult)
}