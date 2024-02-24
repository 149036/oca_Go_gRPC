package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "oca_Go_gRPC/api/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Order struct {
	Id      int
	Content string
}

var (
	order1 = Order{
		Id:      1,
		Content: "Content1",
	}
	order2 = Order{
		Id:      2,
		Content: "Content2",
	}
	orders = []Order{order1, order2}
)

func getOrder(i int32) Order {
	return orders[i-1]
}

type server struct {
	pb.UnimplementedOrderToServer
}

func (s *server) GetOrder(ctx context.Context, in *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order := getOrder(in.Id)

	protoOrder := &pb.Order{
		Id:      int32(order.Id),
		Content: order.Content,
	}
	return &pb.GetOrderResponse{Order: protoOrder}, nil
}

var (
	port = flag.Int("port", 50051, "the port to serve on")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderToServer(s, &server{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
