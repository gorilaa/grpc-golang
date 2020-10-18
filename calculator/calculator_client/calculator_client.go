package main

import (
  "context"
  "fmt"
  "google.golang.org/grpc"
  "grpc-golang/calculator/calculatorpb"
  "io"
  "log"
  "os"
  "strconv"
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

  //Stream
  args, _ := strconv.Atoi(os.Args[1])
  doStreamDecomposition(c, int64(args))
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
	log.Fatalf("Error calling Calculator Rpc: %v", err)
  }

  log.Printf("Response from guets: %v", res.CalculatorResult)
}

func doStreamDecomposition(c calculatorpb.CalculatorServiceClient, args int64)  {
  fmt.Println("Starting to do a Streaming Server Decomposition RPC..")
  req := &calculatorpb.DecompositionRequest{
	Decomposition: &calculatorpb.Decomposition{
	  Parameter: args,
	},
  }
  resDecomposition, err := c.PrimeNumberDecomposition(context.Background(), req)

  if err != nil{
	log.Fatalf("Error calling Stream Server PrimeNumberDecomposition Rpc: %v", err)
  }

  for  {
	msg, err := resDecomposition.Recv()
	if err == io.EOF {
	  // end of the stream
	  break;
	}

	if err !=  nil {
	  log.Fatalf("Error reading stream: %v", err)
	}

	log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetDecompositionResult())
  }
}