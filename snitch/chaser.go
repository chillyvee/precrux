package snitch

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Precrux Structs
type ChaserProfile struct {
	Name        string `json:"name" yaml:"name"`
	Addr        string `json:"addr" yaml:"addr"`
	Certificate string `json:"certificate" yaml:"certificate"`
}

func (cp *ChaserProfile) LoadFromFile(fp string) error {
	log.Printf("Loading %s", fp)
	yamlData, err := os.ReadFile(fp)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return errors.New(fmt.Sprintf("Unable to read file %s; %s", fp, err))
	}
	err = yaml.Unmarshal(yamlData, cp)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return errors.New(fmt.Sprintf("Invalid YAML Synxtax in %s: %s", fp, err))
	}
	return nil
}

func (cp *ChaserProfile) Load() error {
	return cp.LoadFromFile(cp.Filepath())
}

func (cp *ChaserProfile) Filepath() string {
	return fmt.Sprintf("chaser_%s.yaml", cp.Name)
}

func (cp *ChaserProfile) Save() {
	yamlBytes, err := yaml.Marshal(cp)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if err := os.WriteFile(cp.Filepath(), yamlBytes, 0600); err != nil {
		panic(err)
	}
	fmt.Printf("[chaser %s] Profile saved to %s\n", cp.Name, cp.Filepath())
}

func (s *Snitch) AddChaserProfile(chaserName, chaserAddr string) {
	cpFile := fmt.Sprintf("chaser_%s.yaml", chaserName)
	if _, err := os.Stat(cpFile); !errors.Is(err, os.ErrNotExist) {
		err := errors.New(fmt.Sprintf("Chaser profile for %s already exists: %s", chaserName, cpFile))
		panic(err)
	}

	certificate := s.ReadChaserCertificate(chaserName)

	cp := &ChaserProfile{Name: chaserName, Addr: chaserAddr, Certificate: certificate}
	cp.Save()
}

func (s *Snitch) ReadChaserCertificate(chaserName string) string {
	fmt.Println("Copy certificate from chaser and paste below.  Input a blank line to finish.")
	fmt.Println("First Line Should Start with -----BEGIN CERTIFICATE-----")
	fmt.Println("")

	var textBuffer strings.Builder
	textBuffer.Grow(10240)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			break
		}
		textBuffer.WriteString(scanner.Text() + "\n")
	}
	certstring := textBuffer.String()

	certfile := fmt.Sprintf("chaser_%s.crt", chaserName)
	if err := os.WriteFile(certfile, []byte(certstring), 0600); err != nil {
		panic(err)
	}

	fmt.Printf("Chaser %s certificate saved to %s\n", chaserName, certfile)
	return certstring
}
