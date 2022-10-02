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
	//chaserCertBytes []byte

	ChainName string

	ChaserName    string
	ChaserAddress string
	ChaserConn    *grpc.ClientConn
	ChaserRemote  pb.PrecruxClient

	Config SnitchConfig
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

func callWriteFile(client pb.PrecruxClient, filepath string, data []byte, perm os.FileMode, overwrite bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.WriteFile(ctx, &pb.PrecruxWriteFileRequest{Filepath: filepath, Data: data, Perm: uint32(perm), Overwrite: overwrite})
	if err != nil {
		log.Fatalf("client.WriteFile(_) = _, %v: ", err)
	}
	fmt.Println("WriteFile: succes=", resp.Success, " message=", resp.Message)
}

/*
func (s *Snitch) LoadCert() {
	contents, err := os.ReadFile(fmt.Sprintf("chaser_%s.crt", s.ChaserName))
	if err != nil {
		panic(err)
	}
	s.chaserCertBytes = contents

}
*/
func (s *Snitch) Connect() {
	//s.LoadCert()
	cp := &ChaserProfile{Name: s.ChaserName}
	cp.Load()

	fmt.Println("Starting Snitch")
	fmt.Println("chaser-name:", cp.Name)
	fmt.Println("chaser address: ip:port :", cp.Addr)
	fmt.Println("")

	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM([]byte(cp.Certificate)) {
		fmt.Println("credentials: failed to append certificates")
		panic("certificate corrupted")
	}

	//block, _ := pem.Decode(s.chaserCertBytes)
	block, _ := pem.Decode([]byte(cp.Certificate))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	scert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	parts := strings.Split(scert.Subject.String(), "=")
	remoteName := parts[1]

	creds := credentials.NewTLS(&tls.Config{ServerName: remoteName, RootCAs: certpool}) //  okay if ServerName matches.  At least crt content match is required

	// Set up a connection to the server.
	conn, err := grpc.Dial(cp.Addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	s.ChaserConn = conn
	s.ChaserRemote = pb.NewPrecruxClient(conn)
}

func (s *Snitch) SendFile(fp string, co *ChaserConfig) {
	// Make a echo client and send an RPC.
	// callUnaryPrecrux(s.ChaserRemote, "hello world")

	srcDir := s.ChaserDir(co)
	// srcDir := fmt.Sprintf("%s/chaser/%s", s.ChainName, s.ChaserName)
	dstDir := fmt.Sprintf("%s", s.ChainName)

	srcFile := fmt.Sprintf("%s/%s", srcDir, fp)
	dstFile := fmt.Sprintf("%s/%s", dstDir, fp)

	data, err := os.ReadFile(srcFile)
	if err != nil {
		panic(err)
	}
	callWriteFile(s.ChaserRemote, dstFile, data, 0600, true)
}

func (s *Snitch) SendChainToChaser() {
	s.ReadConfig()
	s.Connect()
	defer s.ChaserConn.Close()
	for _, co := range s.Config.Cosigners {
		if co.Name == s.ChaserName {
			s.SendFile("config.yaml", co)
			s.SendFile("share.json", co)
			s.SendFile(s.ChaserStatePvsFile(), co)
			s.SendFile(s.ChaserStateSssFile(), co)
		}
	}
}
