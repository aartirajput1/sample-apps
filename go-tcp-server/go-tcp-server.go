package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

const (
	// LISTENER_NUM is the env var that can be used to define the number of TCP listeners
	// that will be spun up on this server.
	LISTENER_NUM = "LISTENER_NUM"
	// MIN_LISTENERS is the minimum number of TCP listeners available to run at one time.
	MIN_LISTENERS = 1
	// MAX_LISTENERS is the maximum number of TCP listeners available to run at one time.
	MAX_LISTENERS = 10
	// STARTING_PORT is the initial port number that will server a TCP listner.
	STARTING_PORT = 10001
	// MAX_REQUESTS is the number of requests the TCP server will handle before closing
	// the listiner.
	MAX_REQUESTS = 3
)

func main() {
	fmt.Println("Starting TCP listener...")

	listenerNumEnv := os.Getenv(LISTENER_NUM)
	if listenerNumEnv == "" {
		log.Fatalf("Environment variable %q is not defined; exiting", LISTENER_NUM)
	}
	listenerNum, err := strconv.Atoi(listenerNumEnv)
	if err != nil {
		log.Fatalf("%q is not a whole number: %s", LISTENER_NUM, err)
	}
	if listenerNum < MIN_LISTENERS || listenerNum > MAX_LISTENERS {
		log.Fatalf("%q (%d) should be between %d and %d; exiting", LISTENER_NUM, listenerNum, MIN_LISTENERS, MAX_LISTENERS)
	}

	for port := STARTING_PORT; port < STARTING_PORT+listenerNum; port++ {
		go buildAndRunListner(port)
	}

	// We will block here even as TCP listeners begin to close
	select {}
}

// buildAndRunListner will build a TCP listener on the provided port and begin accepting
// connections. The listener will close after serving a given number of connections.
func buildAndRunListner(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer l.Close()
	if err != nil {
		log.Fatalf("Error listening on port %d: %s", port, err)
	}
	fmt.Printf("Server listening on TCP port %d\n", port)

	// Set counter that will terminate the TCP listener when we reach MAX_REQUESTS
	var counter int

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting incoming request: %s", err)
		}

		counter++
		fmt.Printf("Processing request %d on port %d from %s\n", counter, port, conn.RemoteAddr())
		conn.Write([]byte(fmt.Sprintf("Request %d processed", counter)))
		conn.Close()

		if counter >= MAX_REQUESTS {
			l.Close()
			log.Printf("Max requests of %d reached: listener closed on port %d", MAX_REQUESTS, port)
			return
		}
	}
}
