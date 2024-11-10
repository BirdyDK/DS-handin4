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
	commandChannel := make(chan string, 100)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		//log.Println("Scanner ready")
		for scanner.Scan() {
			//log.Println("Scanned: " + scanner.Text())
			commandChannel <- scanner.Text()
			//log.Println("Scan sent to channel")

		}
	}()

	for {
		var command string
		if len(commandChannel) > 0 {
			command = <-commandChannel
		}

		if command == "enter" {
			log.Println("Detected ENTER")
			for {
				if peerNode.HasToken {
					peerNode.EnterCriticalSection()
					break
				}
			}
		}
		if peerNode.HasToken {
			time.Sleep(3 * time.Second)
			peerNode.PassToken()
		}
	}
}
