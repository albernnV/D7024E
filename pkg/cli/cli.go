package cli

import "fmt"
import "strings"

func Cli() {
	for {
		var input string
		fmt.Printf("Command: \n")
		fmt.Scanln(&input)

		switch {
		case strings.Contains(input,"PUT"):
			fmt.Println("put command")
		case strings.Contains(input, "GET"):
			fmt.Println("get command")
		case strings.Contains(input, "PING"):
			fmt.Println("ping command")
		case input == "EXIT":
			fmt.Println("exit command")
		}
	}
}