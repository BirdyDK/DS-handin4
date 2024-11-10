package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BirdyDK/DS-handin4/node"
)

func main() {
	id := flag.String("id", "", "Node ID")
	address := flag.String("address", "", "Node address")
	nextNode := flag.String("nextNode", "", "Address of the next node")
	hasToken := flag.Bool("token", false, "Does this node start with the token")

	flag.Parse()

	peerNode := node.NewNode(*id, *address, *nextNode, *hasToken)

	go node.StartGRPCServer(*address, peerNode)

	// Signal handling for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		peerNode.CloseConnections()
		os.Exit(0)
	}()

	// Command channel
	commandChannel := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			commandChannel <- scanner.Text()
		}
	}()

	for {
		time.Sleep(3 * time.Second)
		var command string
		if len(commandChannel) > 0 {
			log.Println("length of commandChannel: " + string(len(commandChannel)))
			command = <-commandChannel
			log.Println("Command Value: " + command + "; Length of commandChannel: " + string(len(commandChannel)))

		}

		if command == "enter" {
			for {
				if peerNode.HasToken {
					peerNode.EnterCriticalSection()
					break
				}
			}
		}
		peerNode.PassToken()
	}
}
