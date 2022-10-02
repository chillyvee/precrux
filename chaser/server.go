package chaser

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

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

func (s *pbServer) WriteFile(ctx context.Context, req *pb.PrecruxWriteFileRequest) (*pb.PrecruxWriteFileResponse, error) {
	fmt.Printf("WriteFile %+q", req)
	dir := filepath.Dir(req.Filepath)
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating Directory %s\n", dir)
		if err := os.MkdirAll(dir, 0700); err != nil {
			panic(err)
		}
	}

	if !req.Overwrite {
		if _, err := os.Stat(req.Filepath); errors.Is(err, os.ErrNotExist) {
			return &pb.PrecruxWriteFileResponse{Success: false, Message: fmt.Sprintf("Can't overwrite %s", req.Filepath)}, nil
		}
	}

	if err := os.WriteFile(req.Filepath, req.Data, os.FileMode(req.Perm)); err != nil {
		fmt.Printf("Error: %q", err)
		return &pb.PrecruxWriteFileResponse{Success: false, Message: fmt.Sprintf("Write Error %s", err)}, nil
	}
	return &pb.PrecruxWriteFileResponse{Success: true, Message: "hello again"}, nil
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
