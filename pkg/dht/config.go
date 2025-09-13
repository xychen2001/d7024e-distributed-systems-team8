// pkg/dht/config.go
package dht

type Config struct {
	K     int
	Alpha int
	SelfAddr string
	Bootstrap []string // addrs
}