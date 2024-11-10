package node

func (n *Node) DiscoverNodes() {
	for _, addr := range n.nodes {
		client, conn := NewGRPCClient(addr)
		// Store the client and connection
		n.clients[addr] = client
		n.conns[addr] = conn
	}
}
