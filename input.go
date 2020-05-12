package main

import (
	"bufio"
	"fmt"
	"os"
)

func inputFromUser(name string) (inputUser string) {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter " + name + " :")
	if !scanner.Scan() {
		return
	}
	return scanner.Text()
}
