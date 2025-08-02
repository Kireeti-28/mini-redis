package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

enter 'store help' to revist the instructions again.`

const serverBaseURL = "http://localhost:9686/kv"

func main() {
	sayHello()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(">")
	for scanner.Scan() {
		input := scanner.Text()
		inputParts, err := validInput(input)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			switch inputParts[1] {
			case "get":
				storeGet(inputParts[2])
			case "set":
				storeSet(inputParts[2], inputParts[3])
			case "delete":
				storeDelete(inputParts[2])
			case "view":
				storeView()
			case "help":
				fmt.Println(instructionMessage)
			}
		}

		fmt.Printf(">")
	}
}

func sayHello() {
	fmt.Println("Hello, Welcome!👋")
	fmt.Println(instructionMessage)
}

func validInput(input string) ([]string, error) {
	inputParts := strings.Split(strings.TrimSpace(input), " ")
	if len(inputParts) < 2 {
		return nil, fmt.Errorf("invalid input: %v", input)
	}

	if inputParts[0] != "store" {
		return nil, fmt.Errorf("invalid input: %v", input)
	}

	switch inputParts[1] {
	case "get":
		if len(inputParts) != 3 {
			return nil, fmt.Errorf("invalid input: %v", input)
		}
	case "set":
		if len(inputParts) != 4 {
			return nil, fmt.Errorf("invalid input: %v", input)
		}
	case "delete":
		if len(inputParts) != 3 {
			return nil, fmt.Errorf("invalid input: %v", input)
		}
	case "view":
	case "help":
	default:
		return nil, fmt.Errorf("invalid input %v", input)
	}

	return inputParts, nil
}

func storeGet(key string) {
	resp, err := http.Get(serverBaseURL + "/" + key)
	if err != nil {
		fmt.Printf("failed to make request: %v\n", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read resp body %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(string(data))
}

func storeSet(key, value string) {
	resp, err := http.Post(serverBaseURL+"/"+key, "text/plain", strings.NewReader(value))
	if err != nil {
		fmt.Printf("failed to make request: %v\n", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read resp body %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(string(data))
}

func storeDelete(key string) {
	req, err := http.NewRequest(http.MethodDelete, serverBaseURL+"/"+key, nil)
	if err != nil {
		fmt.Printf("failed to make request: %v\n", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed to make reqeust %v\n", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read resp body %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(string(data))
}

func storeView() {
	resp, err := http.Get(serverBaseURL)
	if err != nil {
		fmt.Printf("failed to make request: %v\n", err)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read resp body %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(string(data))
}
