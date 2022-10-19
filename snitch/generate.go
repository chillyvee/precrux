package snitch

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

	"github.com/strangelove-ventures/horcrux/signer"
)

// Private Horcrux Structs from "github.com/strangelove-ventures/horcrux/cmd/horcrux/cmd"
type CosignerPeer struct {
	ShareID int    `json:"share-id" yaml:"share-id"`
	P2PAddr string `json:"p2p-addr" yaml:"p2p-addr"`
}

type CosignerConfig struct {
	Threshold int            `json:"threshold"   yaml:"threshold"`
	Shares    int            `json:"shares" yaml:"shares"`
	P2PListen string         `json:"p2p-listen"  yaml:"p2p-listen"`
	Peers     []CosignerPeer `json:"peers"       yaml:"peers"`
	Timeout   string         `json:"rpc-timeout" yaml:"rpc-timeout"`
}

type ChainNode struct {
	PrivValAddr string `json:"priv-val-addr" yaml:"priv-val-addr"`
}

// Config maps to the on-disk JSON format
type DiskConfig struct {
	PrivValKeyFile *string         `json:"key-file,omitempty" yaml:"key-file,omitempty"`
	ChainID        string          `json:"chain-id" yaml:"chain-id"`
	CosignerConfig *CosignerConfig `json:"cosigner,omitempty" yaml:"cosigner,omitempty"`
	ChainNodes     []ChainNode     `json:"chain-nodes,omitempty" yaml:"chain-nodes,omitempty"`
	DebugAddr      string          `json:"debug-addr,omitempty" yaml:"debug-addr,omitempty"`
}

// Precrux Structs
type ChaserConfig struct {
	ChaserName  string `json:"chaser-name" yaml:"chaser-name"`
	P2PListen   string `json:"p2p-listen"  yaml:"p2p-listen"`
	PrivValAddr string `json:"priv-val-addr" yaml:"priv-val-addr"`
	DebugAddr   string `json:"debug-addr,omitempty" yaml:"debug-addr,omitempty"`
}

type SnitchConfig struct {
	ChainName  string          `json:"chain-name" yaml:"chain-name"`
	ChainID    string          `json:"chain-id" yaml:"chain-id"`
	RpcTimeout string          `json:"rpc-timeout" yaml:"rpc-timeout"`
	Threshold  int64           `json:"threshold"   yaml:"threshold"`
	Shares     int64           `json:"shares" yaml:"shares"`
	Cosigners  []*ChaserConfig `json:"cosigners" yaml:"cosigners"`
}

func (s *Snitch) ReadConfig() {
	yamlData, err := ioutil.ReadFile(fmt.Sprintf("%s/precrux.yaml", s.ChainName))
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlData, &s.Config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func tcpUrl(remote string) string {
	return fmt.Sprintf("tcp://%s", remote)
}
func (s *Snitch) Generate() error {
	s.ReadConfig()
	fmt.Printf("Precrux Config: %+q\n", s.Config)
	if err := s.GatherCertificates(); err != nil {
		return err
	}
	s.SplitValidatorKey()
	s.CreateChaserConfigs()
	return nil
}

func (s *Snitch) GatherCertificates() error {
	// Gather Chaser Certificates
	for _, co := range s.Config.Cosigners {
		cp := ChaserProfile{Name: co.ChaserName}
		fmt.Printf("[chaser %s] Configuration Load\n", co.ChaserName)

		err := cp.Load()
		if err != nil || len(cp.Certificate) < 100 {
			fmt.Printf("[chaser %s] Configuration Error: %s\n", co.ChaserName, err)
			return err
		}
		fmt.Printf("[chaser %s] Configuration OK\n", co.ChaserName)
	}
	return nil
}

func (s *Snitch) SplitValidatorKey() {
	missingShares := 0
	for _, co := range s.Config.Cosigners {
		// Check for existing key shard
		shareFile := s.ChaserShareFile(co)
		if _, err := os.Stat(shareFile); errors.Is(err, os.ErrNotExist) {
			fmt.Printf("Missing Key Share %s\n", shareFile)
			missingShares += 1
		}

	}

	// Check if all shards missing existing shards
	if missingShares == len(s.Config.Cosigners) {

		// Shard validator key
		fmt.Printf("Sharding priv_validator_key.json\n")
		csKeys, err := signer.CreateCosignerSharesFromFile(fmt.Sprintf("%s/priv_validator_key.json", s.ChainName), s.Config.Threshold, s.Config.Shares)
		if err != nil {
			panic(err)
		}

		for idx, shareBytes := range csKeys {
			co := s.Config.Cosigners[idx]
			shareFile := s.ChaserShareFile(co)
			shareDir := filepath.Dir(shareFile)
			if err = os.MkdirAll(filepath.Dir(shareFile), 0700); err != nil {
				fmt.Printf("Failed Creating Directory %s\n", shareDir)
				panic(err)
			}
			if err = signer.WriteCosignerShareFile(shareBytes, shareFile); err != nil {
				fmt.Printf("Failed Writing Share %s\n", shareFile)
				panic(err)
			}
			fmt.Printf("Created Share %s\n", shareFile)
		}
	}
}

func (s *Snitch) CreateChaserConfigs() {

	// Create config
	/*
		chain-id: uni-5
		cosigner:
		  threshold: 2
		  shares: 3
		  p2p-listen: tcp://IP:PORT
		  peers:
		  - share-id: 2
		    p2p-addr: tcp://IP:PORT
		  - share-id: 3
		    p2p-addr: tcp://IP:PORT
		  rpc-timeout: 1500ms
		chain-nodes:
		- priv-val-addr: tcp://SENTRY_IP:PORT
		debug-addr: '[::]:9001'
	*/
	for shareidx, co := range s.Config.Cosigners {
		shareID := shareidx + 1
		fmt.Printf("\n[chaser %s] Creating Config - share-id: %d\n", co.ChaserName, shareID)
		peers := []CosignerPeer{}
		for remoteidx, rco := range s.Config.Cosigners {
			if shareidx == remoteidx {
				continue
			}
			remoteShareID := remoteidx + 1
			peer := CosignerPeer{
				ShareID: remoteShareID,
				P2PAddr: tcpUrl(rco.P2PListen),
			}
			peers = append(peers, peer)
		}
		cc := &CosignerConfig{Threshold: int(s.Config.Threshold), Shares: int(s.Config.Shares), P2PListen: tcpUrl(co.P2PListen), Timeout: s.Config.RpcTimeout, Peers: peers}
		cns := []ChainNode{ChainNode{PrivValAddr: tcpUrl(co.PrivValAddr)}}

		hconfig := &DiskConfig{ChainID: s.Config.ChainID, DebugAddr: co.DebugAddr, CosignerConfig: cc, ChainNodes: cns}
		// fmt.Printf("%+v", hconfig)

		// Output
		ybytes, err := yaml.Marshal(hconfig)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//fmt.Printf("--- config dump:\n%s\n\n", string(ybytes))

		// Write share config
		// configFile := fmt.Sprintf("cosigner/%s/config.yaml", co.ChaserName)
		configFile := s.ChaserFile(co, "config.yaml")
		if err := os.WriteFile(configFile, ybytes, 0600); err != nil {
			fmt.Printf("Unable to Write Cosigner Config %s\n", configFile)
		}
		fmt.Printf("Wrote Cosigner Config %s\n", configFile)
	}

}

func (s *Snitch) ChaserFile(co *ChaserConfig, filepath string) string {
	return fmt.Sprintf("%s/%s", s.ChaserDir(co), filepath)
}
func (s *Snitch) ChaserDir(co *ChaserConfig) string {
	return fmt.Sprintf("%s/cosigner/%s", s.ChainName, co.ChaserName)
}
func (s *Snitch) ChaserStateDir(co *ChaserConfig) string {
	return fmt.Sprintf("%s/cosigner/%s/state", s.ChainName, co.ChaserName)
}
func (s *Snitch) ChaserShareFile(co *ChaserConfig) string {
	return s.ChaserFile(co, "share.json")
}
func (s *Snitch) ChaserStatePvsFile() string {
	return fmt.Sprintf("state/%s_priv_validator_state.json", s.Config.ChainID)
}

func (s *Snitch) ChaserStateSssFile() string {
	return fmt.Sprintf("state/%s_share_sign_state.json", s.Config.ChainID)
}
