package main

import (
	"bufio"
	"fmt"
	"os"
)

const instructionMessage = `Please find the instructions below to make best use:

usage: store command [<args>]

These are the commands that are used:
	get keyname
		-- this command gets the key & value from the store if exists else errors if not exists.
	set keyname value
		-- this command set the key value in store; overwrites if already exists
	delete keyname 
		-- this command deletes the key as keyname from the store; 
	view
		-- this command gives the current view of store

enter 'help' to revist the instructions again.`

func main() {
	sayHello()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(">")
	for scanner.Scan() {
		args := scanner.Text()
		fmt.Println(args)

		fmt.Printf(">")
	}
}

func sayHello() {
	fmt.Println("Hello, Welcome!👋")
	fmt.Println(instructionMessage)
}
