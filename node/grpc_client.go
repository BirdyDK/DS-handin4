package node

import (
	"context"
	"log"

	proto "github.com/BirdyDK/DS-handin4/proto/github.com/BirdyDK/DS-handin4"
	"google.golang.org/grpc"
)

func NewGRPCClient(address string) (proto.NodeClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	client := proto.NewNodeClient(conn)
	return client, conn
}

func (n *Node) RequestRemoteAccess(client proto.NodeClient, nodeID string, timestamp int64) bool {
	req := &proto.AccessRequest{NodeId: nodeID, Timestamp: timestamp}
	resp, err := client.RequestAccess(context.Background(), req)
	if err != nil {
		log.Printf("Error while requesting access: %v", err)
		return false
	}
	return resp.Granted
}

func (n *Node) ReleaseRemoteAccess(client proto.NodeClient, nodeID string) {
	req := &proto.AccessRelease{NodeId: nodeID}
	_, err := client.ReleaseAccess(context.Background(), req)
	if err != nil {
		log.Printf("Error while releasing access: %v", err)
	}
}
