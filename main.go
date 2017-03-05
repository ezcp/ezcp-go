package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.LUTC | log.LstdFlags)

	bitcoin := flag.Bool("bitcoin", false, "get a bitcoin address for registration")
	//passphrase := flag.String("passphrase", "", "encrypt / decrypt file")
	login := flag.String("login", "", "register user and set a durable token")
	help := flag.Bool("help", false, "help")

	flag.Parse()

	if *help || flag.NArg() == 0 {
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
}
