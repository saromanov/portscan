package main

// portscan provides multi scanning of the ports
// each port at the separate goroutine

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	minPortParam = kingpin.Flag("minport", "Scanning from this port").Default("0").Int()
	maxPortParam = kingpin.Flag("maxport", "Scanning to this port").Default("65535").Int()
	hostParam    = kingpin.Flag("host", "target host").Required().String()
)

func scan(host string, minPort, maxPort int) {
	var (
		wg sync.WaitGroup
	)

	wg.Add(maxPort - minPort)
	for i := minPort; i < maxPort; i++ {
		go func(port int) {
			addr := fmt.Sprintf("%s:%d", host, port)
			conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
			if err != nil {
				wg.Done()
				fmt.Printf("Error: %v", err)
				return
			}

			wg.Done()
			err = conn.Close()
			if err != nil {
				fmt.Println("Error: %v", err)
			}

			return
		}(i)
	}

	wg.Wait()
}
func main() {
	kingpin.CommandLine.Help = "Simple tool for scanning ports"
	kingpin.Parse()
	if *minPortParam > *maxPortParam {
		fmt.Println("maxport should be greater then minport")
		os.Exit(1)
	}

	scan(*hostParam, *minPortParam, *maxPortParam)
}
