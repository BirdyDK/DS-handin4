package node

import (
	"context"
	"log"
	"net"

	proto "github.com/BirdyDK/DS-handin4/proto/github.com/BirdyDK/DS-handin4"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedNodeServer
	node *Node
}

func NewGRPCServer(node *Node) *grpc.Server {
	s := grpc.NewServer()
	srv := &server{node: node}
	proto.RegisterNodeServer(s, srv)
	return s
}

func (s *server) ReceiveToken(ctx context.Context, req *proto.TokenMessage) (*proto.TokenResponse, error) {
	log.Printf("Node %s: Received token", s.node.ID)
	s.node.mutex.Lock()
	s.node.HasToken = true
	s.node.mutex.Unlock()

	// Automatically attempt to enter the critical section if token is received
	//go s.node.EnterCriticalSection()

	return &proto.TokenResponse{}, nil
}

func StartGRPCServer(address string, node *Node) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := NewGRPCServer(node)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
