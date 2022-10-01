package chaser

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/chillyvee/precrux/proto/precrux"
	"google.golang.org/grpc/credentials"
)

type pbServer struct {
	pb.UnimplementedPrecruxServer
}

func (s *pbServer) UnaryPrecrux(ctx context.Context, req *pb.PrecruxRequest) (*pb.PrecruxResponse, error) {
	fmt.Printf("UnaryPrecrux")
	return &pb.PrecruxResponse{Message: "hello " + req.Message}, nil
}

func (c *Chaser) Serve() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Listening on port: %d", c.Port)

	cert, err := tls.X509KeyPair(c.certBytes, c.keyBytes)
	if err != nil {
		panic(err)
	}
	creds := credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{cert}})

	s := grpc.NewServer(grpc.Creds(creds))

	// Register PrecruxServer on the server.
	pb.RegisterPrecruxServer(s, &pbServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
