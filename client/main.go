package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"mtg-manager/client/api"
)

func main() {
	var err error

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n===== Login Page =====")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Exit")
		fmt.Print("Select an option: ")

		// Read user input
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			err = api.LoginRequest()
		case "2":
			err = api.LoginRequest()
		case "3":
			fmt.Println("Exiting the program.")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, please try again.")
		}

		if err != nil {
			fmt.Println("Authentication failed:", err)
		} else {
			fmt.Println("Successfully authenticated!")
			showMainMenu()
		}
	}
}

func showMainMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n===== Main Menu =====")
		fmt.Println("1. See library")
		fmt.Println("2. Add card to library")
		fmt.Println("3. Remove card from library")
		fmt.Println("4. Create custom card")
		fmt.Println("5. Logout")
		fmt.Print("Select an option: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Println("print out all cards in the library")
		case "2":
			fmt.Println("menu to add card")
		case "3":
			fmt.Println("menu to remove card")
		case "4":
			fmt.Println("menu to create a custom card")
		case "5":
			fmt.Println("Logging out...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
