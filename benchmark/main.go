package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/http2"
)

// TODO: Add command line support for testing this with many clients
func main() {
	for i := 0; i < 100; i++ {
		go connect(i)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Println("Shutting down clients.")
}

func connect(clientID int) {
	// Create a new HTTP/2 Transport with a custom TLS configuration
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Set this to true only for testing purposes
		},
	}

	// Enable HTTP/2 support
	err := http2.ConfigureTransport(tr)
	if err != nil {
		fmt.Println("Error configuring HTTP/2:", err)
		return
	}

	// Create an HTTP client using the custom transport
	client := &http.Client{
		Transport: tr,
	}

	// Send an HTTP/2 GET request
	resp, err := client.Get("https://localhost:3000/subscribe")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("error: ", err)
		}

		fmt.Printf("client %d message: %s", clientID, string(line[:]))
	}
}
