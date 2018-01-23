package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST   = ""
	CONN_PORT   = "3333"
	CONN_PORT_2 = "4444"
	CONN_TYPE   = "tcp"
)

func main() {
	fmt.Println("Starting tcp go server...")

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

	ll, err2 := net.Listen(CONN_TYPE, ":"+CONN_PORT_2)
	if err2 != nil {
		fmt.Println("Error listening:", err2.Error())
		os.Exit(1)
	}
	defer ll.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT_2)
	counter := 0

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)

		conn1, err1 := ll.Accept()
		counter++
		if err1 != nil {
			fmt.Println("Error accepting: ", err1.Error())
			os.Exit(1)
		}

		time.Sleep(50 * time.Millisecond)
		// Stop after first 3 request on port 2.
		if counter == 5 {
			ll.Close()
		}
		go handleRequest(conn1)
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	fmt.Printf("\nReceived message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

	conn.Write([]byte("Hi there !"))
	conn.Close()
}
