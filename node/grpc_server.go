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

func (s *server) RequestAccess(ctx context.Context, req *proto.AccessRequest) (*proto.AccessRespone, error) {
	granted := s.node.RequestAccess(req.NodeId, req.Timestamp)
	return &proto.AccessRespone{Granted: granted}, nil
}

func (s *server) ReleaseAccess(ctx context.Context, req *proto.AccessRelease) (*proto.AccessRespone, error) {
	s.node.ReleaseAccess(req.NodeId)
	return &proto.AccessRespone{Granted: true}, nil
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
