package node

import (
	"context"
	"log"
	"sync"

	proto "github.com/BirdyDK/DS-handin4/proto/github.com/BirdyDK/DS-handin4"
	"google.golang.org/grpc"
)

type Node struct {
	ID              string
	address         string
	nodes           []string
	clients         map[string]proto.NodeClient
	conns           map[string]*grpc.ClientConn
	mutex           sync.Mutex
	csAccess        chan bool
	requestQueue    map[string]int64
	timestamp       int64
	State           string
	deferredReplies []string
}

const (
	STATE_REST   = "REST"
	STATE_WANTED = "WANTED"
	STATE_HELD   = "HELD"
)

func NewNode(id, address string, nodes []string) *Node {
	return &Node{
		ID:              id,
		address:         address,
		nodes:           nodes,
		clients:         make(map[string]proto.NodeClient),
		conns:           make(map[string]*grpc.ClientConn),
		csAccess:        make(chan bool, 100),
		requestQueue:    make(map[string]int64),
		State:           STATE_REST,
		deferredReplies: []string{},
	}
}

func (n *Node) RequestAccess(nodeID string, timestamp int64) bool {
	log.Println(n.ID + " locking RequestAccess")
	n.mutex.Lock()
	log.Println(n.ID + "  RequestAccess locked")
	defer n.mutex.Unlock()

	// Update the node's logical clock
	n.timestamp = max(n.timestamp, timestamp) + 1
	if n.State == STATE_HELD || (n.State == STATE_WANTED && (n.timestamp < timestamp || (n.timestamp == timestamp && n.ID < nodeID))) {
		// Defer the reply
		n.deferredReplies = append(n.deferredReplies, nodeID)
		log.Println(n.ID + " unlocking RequestAccess")
		return false
	}

	// Update the state to HELD
	n.State = STATE_HELD
	log.Println(n.ID + " unlocking RequestAccess")
	return true
}

func (n *Node) CanEnterCriticalSection() bool {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// Node can enter the Critical Section if it is in the WANTED state
	log.Println(n.ID + " unlocking CanEnterCriticalSection")
	return n.State == STATE_WANTED
}

func (n *Node) ReleaseAccess(nodeID string) {

	n.State = STATE_REST
	for _, node := range n.deferredReplies {
		n.sendReply(node)
	}
	n.deferredReplies = []string{}
	n.csAccess <- true
	log.Println(n.ID + " unlocking ReleaseAccess")
}

func (n *Node) NotifyNodesEntering() {
	log.Printf("attempting to lock ")
	n.mutex.Lock()
	log.Printf("finished locking")
	defer n.mutex.Unlock()

	n.State = STATE_WANTED
	n.timestamp++

	for _, client := range n.clients {
		go func(client proto.NodeClient) {

			req := &proto.AccessRequest{NodeId: n.ID, Timestamp: n.timestamp}
			log.Println(n.ID + " attempting to reqeustaccess")
			_, err := client.RequestAccess(context.Background(), req)
			log.Println(n.ID + " requestaccess finished")
			if err != nil {
				log.Printf("Error while sending access request: %v", err)
			}
		}(client)
	}
	log.Println(n.ID + " unlocking NotifyNodesEntering")
}

func (n *Node) NotifyNodesLeaving() {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.State = STATE_REST

	for _, client := range n.clients {
		go func(client proto.NodeClient) {
			req := &proto.AccessRelease{NodeId: n.ID}
			log.Println(n.ID + " attempting to releaseaccess")
			_, err := client.ReleaseAccess(context.Background(), req)
			log.Println(n.ID + " finished releaseaccess")
			if err != nil {
				log.Printf("Error while sending release notification: %v", err)
			}
		}(client)
	}
	log.Println(n.ID + " unlocking NotifyNodesLeaving")
}

func (n *Node) sendReply(nodeID string) {
	log.Println(n.ID + " attempting to connect client to node")
	client, conn := NewGRPCClient(nodeID)
	defer conn.Close()
	log.Println(n.ID + " succeeded in connection")

	req := &proto.AccessRequest{NodeId: n.ID, Timestamp: n.timestamp}
	log.Println(n.ID + " attempting client.requestAccess")
	_, err := client.RequestAccess(context.Background(), req)
	if err != nil {
		log.Printf("Error while sending reply: %v", err)
	}
	log.Println(n.ID + " succeeded in requestAccess")
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func (n *Node) CloseConnections() {
	// Send ReleaseAccess messages to all other nodes
	for _, client := range n.clients {
		n.ReleaseRemoteAccess(client, n.ID)
	}

	// Close client and server connections
	for _, conn := range n.conns {
		conn.Close()
	}
}
