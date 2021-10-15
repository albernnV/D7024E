package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/albernnV/D7024E/pkg/kademlia"
)

// Inputs will have the structure:
// "put STRING"
// "get STRING"
// "exit"
// "PING" - for displaying and testing the ping command
// "lookup ID"

func Cli(kademliaNode *kademlia.Kademlia) {
	for {
		var inputCommand string
		fmt.Printf("Command: \n")
		fmt.Scanln(&inputCommand)

		switch {
		case strings.Contains(inputCommand, "put "):
			inputString := inputCommand[4:]
			kademliaNode.Store([]byte(inputString))

		case strings.Contains(inputCommand, "get "):
			inputString := inputCommand[4:]
			kademliaNode.LookupData(inputString)

		case inputCommand == "exit":
			fmt.Printf("Terminating node")
			os.Exit(0)

		case strings.Contains(inputCommand, "TEST"):
			kademliaNode.Tes()

		//case strings.Contains(inputCommand, "PING"):
		//kademliaNode.SendPing()

		case strings.Contains(inputCommand, "lookup "):
			inputID := inputCommand[4:]
			ID := kademlia.NewKademliaID(inputID)
			contactToLookup := kademlia.NewContact(ID, "")
			shortlist := kademliaNode.LookupContact(&contactToLookup)
			fmt.Print(shortlist)
		}
	}
}
