package main

import (
  "context"
  "fmt"
  "log"
  "net"

  "google.golang.org/grpc"
  "grpc-golang/calculator/calculatorpb"
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