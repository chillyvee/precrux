package chaser

import "fmt"

type Chaser struct {
	certBytes []byte
	keyBytes  []byte

	Name string
}

func (c *Chaser) Start() {
	fmt.Println("Starting Chaser")
	fmt.Println("chaser-name:", c.Name)
	fmt.Println("")
	c.ProvisionCertificate()
}
