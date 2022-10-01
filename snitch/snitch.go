package snitch

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/chillyvee/precrux/proto/precrux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Snitch struct {
	chaserCertBytes []byte

	ChaserName    string
	ChaserAddress string
	ChaserConn    *grpc.ClientConn
	ChaserRemote  pb.PrecruxClient
}

func callUnaryPrecrux(client pb.PrecruxClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryPrecrux(ctx, &pb.PrecruxRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryPrecrux(_) = _, %v: ", err)
	}
	fmt.Println("UnaryPrecrux: ", resp.Message)
}

func (s *Snitch) LoadCert() {
	contents, err := os.ReadFile(fmt.Sprintf("chaser_%s.crt", s.ChaserName))
	if err != nil {
		panic(err)
	}
	s.chaserCertBytes = contents

}
func (s *Snitch) Connect() {
	s.LoadCert()

	fmt.Println("Starting Snitch")
	fmt.Println("chaser-name:", s.ChaserName)
	fmt.Println("chaser address: ip:port :", s.ChaserAddress)
	fmt.Println("")

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(s.chaserCertBytes) {
		fmt.Println("credentials: failed to append certificates")
		panic("certificate corrupted")
	}

	block, _ := pem.Decode(s.chaserCertBytes)
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	scert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	parts := strings.Split(scert.Subject.String(), "=")
	remoteName := parts[1]

	creds := credentials.NewTLS(&tls.Config{ServerName: remoteName, RootCAs: cp}) //  okay if ServerName matches.  At least crt content match is required

	// Set up a connection to the server.
	conn, err := grpc.Dial(s.ChaserAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.ChaserConn = conn
	s.ChaserRemote = pb.NewPrecruxClient(conn)
}

func (s *Snitch) Query() {
	// Make a echo client and send an RPC.
	callUnaryPrecrux(s.ChaserRemote, "hello world")

}

func (s *Snitch) Start() {
	s.Connect()
	defer s.ChaserConn.Close()
	s.Query()
}
