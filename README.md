# Multi-Client Chat Application

## Overview
This project is a multi-client chat application that demonstrates the use of Remote Procedure Calls (RPC) in a distributed systems environment. It consists of two main components: a server that handles incoming messages and client applications that send messages to the server. Initially developed as a single client-server model, it has been extended to support multiple clients.

## Features
- **Single Client-Server Communication**: Initially supports a one-on-one chat between a client and the server.
- **Multi-Client Support**: Enhanced to allow multiple clients to connect and communicate through the server.
- **RPC Implementation**: Uses RPC for sending and receiving messages between clients and the server.
- **Real-time Messaging**: Clients can send and receive messages in real-time.
- **Scalable Architecture**: Designed to be scalable to accommodate more clients with minimal changes.

## Technologies Used
- Golang
- Standard Golang library: `net/rpc``

## Getting Started

### Prerequisites
- Lorem ipsum

### Installation
1. Clone the repository:
   ```bash
   git clone [repository URL]
   ```
2. Navigate to the project directory:
   ```bash
   cd [project directory]
   ```
3. Install dependencies:
   ```bash
   [commands to install dependencies]
   ```

### Running the Application
1. Start the server:
   ```bash
   [command to start the server]
   ```
2. In a new terminal, start a client:
   ```bash
   [command to start a client]
   ```
3. Repeat step 2 to open multiple clients.

## Usage
- Lorem ipsum

## Design
### Base
All nodes will essentially be running the same logic but with different paths depending on their role. The basic logic will establish the listening RPC logic. When a node receives a message, it will store the user's IP in order to keep track of all clients participating.

### Client
The client nodes will spin off the base by reacting to messages on the handler by showing them to the user. They will also establish a connection to the server node and begin relaying messages from the user.

#### Server
The server node differs from the base by storing messages sent to it from clients, storing the IP and sending the message to all clients in the send list. If a client does not ack the message, it will be removed from the stored list