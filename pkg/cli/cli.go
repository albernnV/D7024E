package cli

import (
	"bufio"
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
	reader := bufio.NewReader(os.Stdin)
	for {
		cli := bufio.NewScanner(reader)
		fmt.Printf("Command: \n")
		cli.Scan()
		inputCommand := cli.Text()

		switch {
		case strings.Contains(inputCommand, "put"):
			inputString := inputCommand[4:]
			kademliaNode.Store([]byte(inputString))

		case strings.Contains(inputCommand, "get "):
			inputString := inputCommand[4:]
			kademliaNode.LookupData(inputString)

		case inputCommand == "exit":
			fmt.Printf("Terminating node")
			os.Exit(0)

		case strings.Contains(inputCommand, "lookup "):
			inputID := inputCommand[7:]
			ID := kademlia.NewKademliaID(inputID)
			contactToLookup := kademlia.NewContact(ID, "")
			shortlist := kademliaNode.LookupContact(&contactToLookup)
			fmt.Print(shortlist)
		}
	}
}
