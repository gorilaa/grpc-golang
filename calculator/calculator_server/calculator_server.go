package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"grpc-golang/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct {}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {

  fmt.Println("Sum number", req)
  parameterOne := req.GetCalculator().GetParameterOne()
  parameterTwo := req.GetCalculator().GetParameterTwo()
  result := parameterOne + parameterTwo

  res := &calculatorpb.CalculatorResponse{
	CalculatorResult: result,
  }

  return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.DecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
  fmt.Println("Start Streaming PrimeNumberDecomposition",  req)
  parameter := req.GetDecomposition().GetParameter()
  divisor := int64(2)

  for parameter > 1 {
    if parameter%divisor == 0 {
      err := stream.Send(&calculatorpb.DecompositionResponse{
        DecompositionResult: divisor,
      })

      if err == nil {
        parameter = parameter / divisor
      }
    } else {
      divisor++
      fmt.Println("Divisor has increase to ", divisor)
    }
  }

  return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
  fmt.Printf("ComputeAverage function was invoked")
  
  sum := int32(0)
  count := 0

  for  {
		req, err := stream.Recv()
		if err == io.EOF {
      // finished reading stream
      average := float64(sum) / float64(count)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				ComputeAverageResult: average,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading ComputeAverage stream: %v", err)
		}

		sum += req.GetAverage()
    count++
    
    fmt.Println("Sum: ", sum)
    fmt.Println("Count: ",count)
  }
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
  fmt.Printf("Find Maximum RPC")
  
  maximum := int32(0)
  
  for {
    req, err := stream.Recv()
    if err == io.EOF {
      return nil
    }
    
    if err != nil {
      log.Fatalf("Error reading client stream : %v", err)
      return err
    }
    number := req.GetNumber()
    
    if number > maximum {
      maximum = number
      sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
        Maximum: maximum,
      })
      
      if sendErr != nil {
        log.Fatalf("Error while sending data to client: %v", err)
        return err
      }
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
  calculatorpb.RegisterCalculatorServiceServer(s, &server{})

  if err := s.Serve(lis); err != nil {
	log.Fatalf("Failed to load server: %v", err)
  }
}