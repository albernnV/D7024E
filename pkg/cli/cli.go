package cli

import "fmt"

func StartCli() {
    fmt.Print("Enter text: ")
    var input string
    fmt.Scanln(&input)
    fmt.Print(input)
}