package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"basic-rpc-chat/client"
	"basic-rpc-chat/server"
	"basic-rpc-chat/shared"
)

const (
	clientName = "client"
	serverName = "server"
)

var (
	serverAddress = "0.0.0.0:5083"
	localAddress  = "0.0.0.0"
	ActiveUsers   map[string]bool
	nodeType      string
	userName      string
	port          string
)

func main() {
	// Config setup
	var err error
	initFlags()
	flag.Parse()
	file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Failed reading the config file: %v\n", err)
	}
	var config map[string]string
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal: %s\n", err)
	}
	serverAddress = config["serverIP"]
	if configClientIP, ok := config["clientIP"]; !ok {
		getOutboundIP()
		localAddress = localAddress + ":" + port
	} else {
		localAddress = configClientIP + ":" + port
	}

	// Prep client or server listen handlers
	fmt.Printf("NodeType: %s\n", nodeType)
	if nodeType == clientName {
		err = client.StartClient(localAddress, serverAddress, userName)
	} else if nodeType == serverName {
		err = server.StartServer(serverAddress)
		localAddress = serverAddress
	}
	if err != nil {
		log.Fatalf("Failed to start node: %s\n", err)
	}

	l, err := net.Listen("tcp", localAddress)
	if err != nil {
		log.Fatal("Failed to listen on address: ", err)
	}

	shared.ConsoleMutex.Lock()
	fmt.Printf("Starting %s on %s\n", nodeType, localAddress)
	shared.ConsoleMutex.Unlock()
	err = http.Serve(l, nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %s\n", err)
	}

}

func initFlags() {
	flag.StringVar(&nodeType, "type", "client", "Set the node to either be a client or server")
	flag.StringVar(&port, "port", "5083", "The port for listening to server messages")
	flag.StringVar(&userName, "name", "anonymous", "Set your username")
}

// Get preferred outbound ip of this machine
func getOutboundIP() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress = string(conn.LocalAddr().(*net.UDPAddr).IP.String())
}
