package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/BirdyDK/DS-handin4/node"
)

func main() {
	id := flag.String("id", "", "Node ID")
	address := flag.String("address", "", "Node address")
	nodes := flag.String("nodes", "", "Comma-separated list of other nodes")

	flag.Parse()

	// Ensure the addresses include ports
	if !strings.Contains(*address, ":") {
		log.Fatalf("Invalid address: %s. Must include port.", *address)
	}

	peerNode := node.NewNode(*id, *address, strings.Split(*nodes, ","))

	go node.StartGRPCServer(*address, peerNode)
	peerNode.DiscoverNodes()

	// Signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		peerNode.CloseConnections()
		os.Exit(0)
	}()

	// Informational log
	log.Println("Node", peerNode.ID, "started and ready to accept commands.")
	log.Println("To try to enter the Critical Section, type 'enter' and press Enter.")

	// Read user input from the command line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()
		if command == "enter" {
			go attemptCriticalSection(peerNode)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func attemptCriticalSection(peerNode *node.Node) {
	log.Println("Node", peerNode.ID, "is attempting to enter the Critical Section...")
	peerNode.NotifyNodesEntering()

	// Log the state after notifying nodes
	log.Println("Node state after NotifyNodesEntering:", peerNode.State)

	time.Sleep(2 * time.Second) // Wait for replies

	// Log the state before checking if it can enter the Critical Section
	log.Println("Node state before CanEnterCriticalSection check:", peerNode.State)

	if peerNode.CanEnterCriticalSection() {
		log.Println("Node", peerNode.ID, "entered the Critical Section")
		time.Sleep(1 * time.Second) // Simulate Critical Section
		peerNode.NotifyNodesLeaving()
		log.Println("Node", peerNode.ID, "left the Critical Section")
	} else {
		log.Println("Node", peerNode.ID, "could not enter the Critical Section at this time.")
	}
}
