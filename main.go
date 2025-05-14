// file: main.go
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// buildRIPRequest constructs a RIP v2 request packet.
// It sets Command=1 (Request), Version=2, and includes a single
// entry (AFI=0, Metric=16) to request the entire routing table.
func buildRIPRequest() []byte {
	pkt := make([]byte, 4+20)
	pkt[0] = 1 // Request
	pkt[1] = 2 // Version
	// Metric at last 4 bytes of the entry
	binary.BigEndian.PutUint32(pkt[4+16:], 16)
	return pkt
}

// parseRIPResponse parses RIP entries from data and prints each one.
func parseRIPResponse(data []byte, addr *net.UDPAddr) {
	if len(data) < 4 {
		fmt.Printf("Received too-short packet from %v\n", addr)
		return
	}
	cmd := data[0]
	version := data[1]
	if cmd != 2 {
		fmt.Printf("Ignoring non-response packet (cmd=%d) from %v\n", cmd, addr)
		return
	}
	fmt.Printf("RIP v%d response from %v:\n", version, addr)

	entries := (len(data) - 4) / 20
	for i := 0; i < entries; i++ {
		base := 4 + i*20
		afi := binary.BigEndian.Uint16(data[base : base+2])
		tag := binary.BigEndian.Uint16(data[base+2 : base+4])
		dst := net.IP(data[base+4 : base+8])
		mask := net.IP(data[base+8 : base+12])
		nxt := net.IP(data[base+12 : base+16])
		metric := binary.BigEndian.Uint32(data[base+16 : base+20])
		fmt.Printf("  Entry %d: AFI=%d, Tag=%d, Dest=%s, Mask=%s, NextHop=%s, Metric=%d\n",
			i+1, afi, tag, dst, mask, nxt, metric)
	}
}

func main() {
	target := flag.String("target", "", "Destination IP or multicast address to send RIP request")
	port := flag.Int("port", 520, "Destination UDP port (default 520)")
	flag.Parse()

	if *target == "" {
		fmt.Println("Error: -target is required")
		flag.Usage()
		os.Exit(1)
	}

	// Listen on UDP port 520 (requires root privileges)
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: 520})
	if err != nil {
		fmt.Printf("Error binding to port 520: %v\n", err)
		return
	}
	defer conn.Close()

	dst := &net.UDPAddr{
		IP:   net.ParseIP(*target),
		Port: *port,
	}

	req := buildRIPRequest()
	if _, err := conn.WriteToUDP(req, dst); err != nil {
		fmt.Printf("Failed to send RIP request to %v: %v\n", dst, err)
		return
	}
	fmt.Printf("Sent RIP request to %v, waiting for responses...\n", dst)

	if err = conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return
	}
	buf := make([]byte, 512)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			break // timeout or other error
		}
		parseRIPResponse(buf[:n], addr)
	}

	fmt.Println("Done.")
}
