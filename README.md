Node ID (--id): A unique identifier for the node. This helps distinguish between different nodes in the system.

Node Address (--address): The IP address and port on which the node listens for incoming gRPC requests. This allows other nodes to communicate with this node.

Next Node (--nextNode): The IP address and port on which the node will pass the token to.

Token (--token): It is false by default, but you have to choose 1 of the nodes in the network to have the token(set token to true) else the system won't work.

**Example of how to run the program with 3 nodes:**

go run main.go -id=node1 -address="localhost:5001" -nextNode="localhost:5002" -token=true

go run main.go -id=node2 -address="localhost:5002" -nextNode="localhost:5003"

go run main.go -id=node3 -address="localhost:5003" -nextNode="localhost:5001"
