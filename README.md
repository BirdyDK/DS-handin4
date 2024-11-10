Node ID (--id): A unique identifier for the node. This helps distinguish between different nodes in the system.

Node Address (--address): The IP address and port on which the node listens for incoming gRPC requests. This allows other nodes to communicate with this node.

Nodes (--nodes): A comma-separated list of addresses of other nodes in the system. This list is used for the initial discovery of other nodes so the node knows where to send messages for coordination.

Example of how to run the program with 3 nodes:
go run main.go --id=node1 --address="localhost:5001" --nodes="localhost:5002,localhost:5003"  
go run main.go --id=node2 --address="localhost:5002" --nodes="localhost:5003,localhost:5001"  
go run main.go --id=node3 --address="localhost:5003" --nodes="localhost:5001,localhost:5002"
