package chaser

import "fmt"

type Chaser struct {
	certBytes []byte
	keyBytes  []byte

	Name string
	Port int
}

func (c *Chaser) Start() {
	fmt.Println("Starting Chaser")
	fmt.Println("chaser-name:", c.Name)
	fmt.Println("port:", c.Port)
	fmt.Println("")
	c.ProvisionCertificate()
	c.Serve()
}
