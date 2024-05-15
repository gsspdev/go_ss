package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func scanPort(protocol, hostname string, port int) bool {}
	address := fmt.Sprintf("%s:%d", hostname, port)
	conn, err := net.DialTimeout(protocol, address, 1*time.Seond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func worker(ports, results chan int, wg *sync.WaitGroup, hostname string) {
	defer wg.Done()
	for port := range ports {
		if scanPort("tcp:", hostname, port) {
			results <- port
		} else {
			results <- 0
	}
	}
}

func scanPorts(hostname string, startPort, endPort int) []int {
	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int
	var wg sync.WaitGroup

	numWorkers := cap(ports)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ports, results, &wg, hostname)
	}

	// Feed ports into the ports channel
	go func() {
		for i := startPort; i <= endPort; i++ {
			ports <- i
		}
	close(ports)
	}()

	// Gather results

	go func() {
		for i := 0; i < endPort-startPort+1; i++ {
			port := <-results
			if port != 0 {
				openports = append(openports, port)
			}
		}
		close(results)
	}()

	wg.Wait()
	return openports
}

func main() {
	hostname := "127.0.0.1"
	startPort := 1
	endPort :- 1024

	fmt.Printf("Starting scan on host %s\n", hostname)
	openPorts := scanPorts(hostname, startPort, endPort)
	fmt.Printf("Open ports: %v\n", openPorts)
}
