package main

import (
	"context"
	"fmt"
	"grpc-golang/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

  if err != nil {
    log.Fatalf("Error calling Greet Rpc: %v", err)
  }

  log.Printf("Response from guets: %v", res.CalculatorResult)

  // doSumUnary(c)

  //Stream
  // args, _ := strconv.Atoi(os.Args[1])
  // doStreamDecomposition(c, int64(args))
  
  // doStreamComputeAverage(c)
  // doFindMaximumStreaming(c)
  doErrorUnary(c)
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

func doStreamComputeAverage(c calculatorpb.CalculatorServiceClient) {
  fmt.Println("Starting ComputeAverage client streaming RPC ...")
  
  stream, err := c.ComputeAverage(context.Background())
  if err != nil {
	  log.Fatalf("Error while calling long greet: %v", err)
  }
  
  numbers := []int32{3, 5, 9, 54, 23}
  
  for _, number := range numbers {
    sendErr := stream.Send(&calculatorpb.ComputeAverageRequest{
      Average: number,
    })
		
		if sendErr != nil {
			log.Fatalf("Error sending data: %v", err)
		}
  }
  
  res, err := stream.CloseAndRecv()
  if err != nil {
		log.Fatalf("Error while recieving response from Compute Average: %v", err)
	}
	fmt.Printf("ComputeAverage Response : %v\n", res.GetComputeAverageResult())
}

func doFindMaximumStreaming(c calculatorpb.CalculatorServiceClient) {
  fmt.Println("Starting FindMaximum client streaming RPC ...")
  
  stream, err := c.FindMaximum(context.Background())
  if err != nil {
    log.Fatalf("Error while calling long greet: %v", err)
  }
  
  waitc := make(chan struct{})
  // send goroutine
  go func() {
    numbers := []int32{4, 7, 2, 19, 4, 6, 32}
    for _, number := range numbers {
      fmt.Printf("Sending number : %v\n", number)
      sendErr := stream.Send(&calculatorpb.FindMaximumRequest{
        Number: number,
      })
      
      if sendErr != nil {
        log.Fatalf("Error sending data: %v", err)
      }
      time.Sleep(1000 * time.Millisecond)
    }
    
    stream.CloseSend()
  }()
  
  // recieve goroutine
  go func() {
    for {
      res, err := stream.Recv()
      if err == io.EOF {
        break
      }
      
      if err != nil {
        log.Fatalf("Error while reading server stream: %v", err)
        break
      }
      
      maximum := res.GetMaximum()
      fmt.Printf("Recieve a new maximum of ...: %v\n", maximum)
    }
    close(waitc)
  }()
  <-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
  fmt.Println("Starting Square Root Unary RPC ...")
  // correct calling
  doErrorCall(c, 10)
  
  // error calling
  doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
  res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
  if err != nil {
    resErr, ok := status.FromError(err)
    if ok {
      // actual error from grpc
      fmt.Printf("Error message from server: %v\n", resErr.Message())
      fmt.Println(resErr.Code())
      
      if resErr.Code() == codes.InvalidArgument {
        fmt.Println("We Probably send a negative number")
      }
    } else {
      log.Fatalf("Big error calling SquareRoot : %v", err)
    }
  }
  
  fmt.Printf("Result of square root of  %v : %v\n ", n, res.GetNumberRoot())
}