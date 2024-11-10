package node

import (
	"context"
	"log"
	"sync"
	"time"

	proto "github.com/BirdyDK/DS-handin4/proto/github.com/BirdyDK/DS-handin4"
	"google.golang.org/grpc"
)

type Node struct {
	ID       string
	address  string
	clients  map[string]proto.NodeClient
	conns    map[string]*grpc.ClientConn
	mutex    sync.Mutex
	HasToken bool
	nextNode string
}

func NewNode(id, address, nextNode string, hasToken bool) *Node {
	return &Node{
		ID:       id,
		address:  address,
		clients:  make(map[string]proto.NodeClient),
		conns:    make(map[string]*grpc.ClientConn),
		HasToken: hasToken,
		nextNode: nextNode,
	}
}

// New function to create a gRPC client connection
func NewGRPCClient(address string) (proto.NodeClient, *grpc.ClientConn) {
	//log.Printf("Attempting to connect to %s", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect to %s: %v", address, err)
		return nil, nil
	}
	client := proto.NewNodeClient(conn)
	/*if client == nil {
		log.Printf("NewNodeClient returned nil for address %s", address)
	} else {
		log.Printf("NewNodeClient successfully created for address %s", address)
	}*/
	return client, conn
}

// Function to pass the token to the next node
func (n *Node) PassToken() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.nextNode == "" || !n.HasToken {
		return
	}

	log.Printf("Node %s: Passing token to %s", n.ID, n.nextNode)
	client, exists := n.clients[n.nextNode]
	if !exists {
		//log.Printf("Node %s: Client does not exist for next node %s", n.ID, n.nextNode)
		// If client doesn't exist, create a new one
		client, conn := NewGRPCClient(n.nextNode)
		if client != nil {
			n.clients[n.nextNode] = client
			n.conns[n.nextNode] = conn
			//log.Printf("Node %s: Stored new client for next node %s", n.ID, n.nextNode)
		} else {
			//log.Printf("Node %s: Could not connect to next node %s", n.ID, n.nextNode)
			return
		}
	}

	// Check if client is nil before calling its methods
	client, exists = n.clients[n.nextNode] // Re-fetch the client to see if it was stored correctly
	if client == nil {
		//log.Printf("Node %s: Re-fetched client for next node %s is nil (exists: %v)", n.ID, n.nextNode, exists)
		return
	}

	// Try to send the token
	_, err := client.ReceiveToken(context.Background(), &proto.TokenMessage{})
	if err != nil {
		//log.Printf("Node %s: Error sending token to %s: %v", n.ID, n.nextNode, err)
	} else {
		n.HasToken = false
	}
}

// Function to receive the token
func (n *Node) ReceiveToken(ctx context.Context, token *proto.TokenMessage) (*proto.TokenResponse, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	//log.Printf("Node %s: Received token", n.ID)
	n.HasToken = true
	//go n.EnterCriticalSection()
	return &proto.TokenResponse{}, nil
}

// Attempt to enter the critical section if the node has the token
func (n *Node) EnterCriticalSection() {
	n.mutex.Lock()
	if !n.HasToken {
		n.mutex.Unlock()
		return
	}
	n.mutex.Unlock()

	log.Printf("Node %s: Entering Critical Section", n.ID)
	time.Sleep(1 * time.Second) // Simulate critical section work
	log.Printf("Node %s: Leaving Critical Section", n.ID)

	// Pass the token to the next node
	//n.PassToken()
}

func (n *Node) CloseConnections() {
	for _, conn := range n.conns {
		conn.Close()
	}
}
