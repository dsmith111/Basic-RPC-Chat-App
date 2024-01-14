package server

import (
	"basic-rpc-chat/shared"
	"fmt"
	"net/rpc"
)

type ServerMessageController struct {
	connectedUsers map[string]bool
	vectorClock    map[string]int
}

func (message *ServerMessageController) Send(messageData shared.Message, ack *shared.Message) error {
	reply := &shared.Message{}
	atLeastOne := false

	// This server controller should send messages to all connected users
	if messageData.Ack == 0 {
		if _, ok := message.connectedUsers[messageData.IpAddress]; !ok {
			fmt.Printf("Adding user with IP %s\n", messageData.IpAddress)
			message.connectedUsers[messageData.IpAddress] = true
			message.vectorClock[messageData.IpAddress] = 0
		}

		if messageData.MessagesSent < message.vectorClock[messageData.IpAddress] {
			ack.Ack = 1
			return nil
		} else if messageData.MessagesSent > message.vectorClock[messageData.IpAddress] {
			message.vectorClock[messageData.IpAddress] = messageData.MessagesSent
		}

		message.vectorClock[messageData.IpAddress] += 1

		fmt.Printf("Delivering message: %s\n", messageData.Data)

		// This is a request, disseminate
		for user, active := range message.connectedUsers {
			if !active {
				continue
			}
			// Dial
			client, err := rpc.DialHTTP("tcp", user)
			if err != nil {
				message.connectedUsers[user] = false
				fmt.Printf("Failed to dial user node: %s\n", err)
				continue
			}

			// RPC
			err = client.Call("ClientMessageController.Send", messageData, reply)
			if err != nil {
				fmt.Printf("Failed to RPC user with error: %s", err)
			}
			atLeastOne = true
		}

		if !atLeastOne {
			err := fmt.Errorf("All calls failed to users")
			message.vectorClock[messageData.Data] -= 1
			return err
		}
	}

	// We will implement ack handling a bit later
	ack.Ack = 1
	return nil
}

/*This method will start an RPC server that will push/pull messages to/from all seen clients*/
/*We'll want to have two routines, a sender and a listener
 */
func StartServer(localAddress string) error {
	serverMessageController := new(ServerMessageController)
	serverMessageController.connectedUsers = make(map[string]bool)
	serverMessageController.vectorClock = make(map[string]int)
	err := rpc.Register(serverMessageController)
	if err != nil {
		fmt.Printf("Error while registering server controller: %s", err)
		return err
	}
	rpc.HandleHTTP()

	return nil
}
