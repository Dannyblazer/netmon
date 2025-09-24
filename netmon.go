package main

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func ping(host string) (time.Duration, error) {
	c, err := net.Dial("ip4:icmp", host)
	if err != nil {
		return 0, err
	}
	defer c.Close()

	// Create ICMP Echo Request
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{ID: 1, Seq: 1},
	}
	bytes, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}

	// Measure Time
	start := time.Now()
	_, err = c.Write(bytes)
	if err != nil {
		return 0, err
	}

	// Wait for reply
	reply := make([]byte, 1500)
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err = c.Read(reply)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func monitorServer(host string) {
	for {
		duration, err := ping(host)
		if err != nil {
			fmt.Printf("Ping %s failed: %v\n", host, err)
		} else {
			fmt.Printf("Ping %s: %v\n", host, duration)
		}
		time.Sleep(5 * time.Second)
	}

}

func main() {
	fmt.Println("Starting network monitor...")
	go monitorServer("8.8.8.8")
	select {}
}
