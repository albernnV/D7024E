package cli

import ("fmt"
		"strings"
		"os"
)

// Inputs will have the structure:
// "put STRING"
// "get STRING"
// "exit"
// "PING" - for displaying and testing the ping command
// "lookup ID"

func Cli(net *Network) {
	for {
		var input string
		fmt.Printf("Command: \n")
		fmt.Scanln(&inputCommand)

		switch {
		case strings.Contains(inputCommand,"put "):
			inputString := inputCommand[4:]

		case strings.Contains(inputCommand, "get "):
			inputString := inputCommand[4:]
			ID := NewKademliaID(inputID)

			SendFindDataMessage()

		case inputCommand == "exit":
			fmt.Printf("Terminating node")
			os.Exit(0)
			
		case strings.Contains(inputCommand, "PING"):
			net.SendPingMessage()

		case strings.Contains(inputCommand, "lookup "):
			inputID := inputCommand[4:]
			ID := NewKademliaID(inputID)
			SendFindContactMessage()
		}
	}
}
