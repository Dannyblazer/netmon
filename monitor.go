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

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{ID: 1, Seq: 1},
	}
	bytes, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()
	_, err = c.Write(bytes)
	if err != nil {
		return 0, err
	}

	reply := make([]byte, 1500)
	c.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, err = c.Read(reply)
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}

func monitorWithAlerts(host string) {
	failCount, highLatencyCount := 0, 0
	for {
		duration, err := ping(host)
		if err != nil {
			failCount++
			fmt.Printf("Ping %s failed: %v\n", host, err)
			if failCount >= 3 {
				fmt.Printf("ALERT: Server %s is down after %d failures!\n", host, failCount)
			}
		} else {
			fmt.Printf("Ping %s: %v\n", host, duration)
			failCount = 0 // Reset on success
		}
		if duration > 200*time.Millisecond {
			highLatencyCount++
			if highLatencyCount >= 4 {
				fmt.Printf("WARNING: High latency to %s: %v\n", host, duration)
				// Or Send email alert
			}
		} else {
			highLatencyCount = 0 // Reset latency to normal
		}
		time.Sleep(5 * time.Second)
	}
}

func main() {
	fmt.Println("Starting network monitor...")
	go monitorWithAlerts("198.54.115.193") // IP to ping and monitor
	select {}
}
