package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LUTC | log.LstdFlags)

	bitcoin := flag.Bool("bitcoin", false, "get a bitcoin address for registration")
	//passphrase := flag.String("x", "", "encrypt / decrypt file")
	help := flag.Bool("help", false, "output usage information")
	//version := flag.Bool("version", false, "output the version number")
	passphrase := flag.String("x", "", "encrypt / decrypt file")
	login := flag.String("login", "", "register user and set a durable token")

	flag.Parse()

	if *help {
		showHelp()
		return
	}
	if *bitcoin {
		address, err := getBitcoinAddress()
		if err != nil {
			panic(err)
		}
		fmt.Println("Please make your 0.01 BTC paiement to: " + address)
		return
	}
	if *login != "" {
		token, err := getToken(*login)
		if err != nil {
			panic(err)
		}
		fmt.Println("Here's your permanent token: " + token)
		return
	}

	switch flag.NArg() {
	case 0:
		if stdinTerminal() && stdoutTerminal() {
			showHelp()
			return
		}
		if stdinTerminal() {
			token, err := getDurableToken()
			if err != nil {
				fmt.Print("You have to login first!")
				return
			}
			err = download(*passphrase, os.Stdout, token)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if stdoutTerminal() {
			token, err := getDurableToken()
			if err != nil {
				fmt.Print("You have to login first!")
				return
			}
			err = upload(*passphrase, os.Stdin, token)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
	}
}
