package client

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"time"

	"basic-rpc-chat/shared"
)

const (
	maxRetries = 3
	retryDelay = time.Second * 3
)

var (
	connectedServer *rpc.Client
	hasConnected    = false
	address         string
	messagesSent    = 0
	clientName      string
)

type ClientMessageController struct {
	connectedUsers map[string]bool
}

func (controller *ClientMessageController) Send(messageData shared.Message, ack *shared.Message) error {
	// This is an RPC for a different node. If it's being called, then we should deliver this message
	// locally
	ackMessage := shared.Message{
		IpAddress: address,
		Ack:       1,
	}

	*ack = ackMessage

	if messageData.IpAddress != address {
		displayMessage(&messageData)
	}

	return nil
}

func StartClient(localAddress string, serverAddress string, userName string) error {
	address = localAddress
	clientName = userName

	// Expose controller on server
	clientMessageController := new(ClientMessageController)
	err := rpc.Register(clientMessageController)

	if err != nil {
		fmt.Printf("Error while registering client controller: %s\n", err)
		return err
	}

	rpc.HandleHTTP()

	// Begin dial with retry to the server
	dialServer(serverAddress)

	// Begin user message handling
	go handleUserMessaging()

	return nil
}

func handleUserMessaging() {
	// In a loop get user inputs and send them
	for {
		// First make sure we're connected
		if !hasConnected {
			fmt.Println("Waiting for connection...")
			time.Sleep(retryDelay * 2)
			continue
		}
		scanner := bufio.NewScanner(os.Stdin)
		shared.ConsoleMutex.Lock()
		fmt.Print("User> ")
		scanner.Scan()
		input := scanner.Text()
		shared.ConsoleMutex.Unlock()
		if connectedServer != nil && len(input) > 0 {
			msg := buildMessage(input)
			sendMessageWithRetry(msg)
		}

	}
}

func sendMessageWithRetry(msg *shared.Message) error {
	reply := &shared.Message{}
	messagesSent += 1

	for i := 0; i < maxRetries; i++ {
		err := connectedServer.Call("ServerMessageController.Send", *msg, reply)
		if err != nil {
			fmt.Printf("Failed to RPC server with error: %s\n", err)
			time.Sleep(retryDelay)
			continue
		} else if reply.Ack != 1 {
			fmt.Println("Server failed to ack message, resending")
			time.Sleep(retryDelay)
			continue
		}
	}
	err := fmt.Errorf("Failed to RPC server with retries\n")
	return err
}

func buildMessage(inputText string) *shared.Message {
	message := &shared.Message{
		Data:         inputText,
		User:         clientName,
		IpAddress:    address,
		Ack:          0,
		MessagesSent: messagesSent,
	}

	return message
}

func displayMessage(receivedMessage *shared.Message) {
	fmt.Printf("\n%s> %s\n", receivedMessage.User, receivedMessage.Data)
	fmt.Print("User> ")
}

func dialServer(serverAddress string) error {
	for i := 0; i < maxRetries; i++ {
		server, err := rpc.DialHTTP("tcp", serverAddress)
		if err != nil {
			hasConnected = false
			fmt.Printf("Failed to dial server with error: %s\n", err)
			time.Sleep(retryDelay)
			continue
		}
		connectedServer = server
		shared.ConsoleMutex.Lock()
		hasConnected = true
		fmt.Printf("Succesfully dialed server on %s\n", serverAddress)
		shared.ConsoleMutex.Unlock()
		return nil
	}
	err := fmt.Errorf("Failed to dial server with retries\n")
	return err
}
