package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	Version = "1.0.0"
)

func main() {
	log.SetFlags(log.LUTC | log.LstdFlags)

	bitcoin := flag.Bool("b", false, "get a bitcoin address for registration")
	help := flag.Bool("help", false, "output usage information")
	version := flag.Bool("version", false, "output the version number")
	passphrase := flag.String("x", "", "encrypt / decrypt file")
	login := flag.String("l", "", "register user and set a durable token")

	flag.Parse()

	if *version {
		fmt.Println("EZCP-go v" + Version)
		return
	}
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
	case 1:
		token, err := getDurableToken()
		if err == nil && stdinTerminal() && stdoutTerminal() && !isSHA1Token(flag.Arg(0)) {
			fpath := flag.Arg(0)
			fileinfo, err := os.Stat(fpath)
			if err == nil && !fileinfo.IsDir() {
				// upload mode
				file, err := os.Open(fpath)
				if err != nil {
					log.Fatal(err)
				}
				err = upload(*passphrase, file, token)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			// download mode
			file, err := os.Create(fpath)
			if err != nil {
				log.Fatal(err)
			}
			err = download(*passphrase, file, token)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		token = flag.Arg(0)
		if !isSHA1Token(token) {
			log.Fatalf("%s is not a valid token", token)
		}
		if stdinTerminal() {
			err := download(*passphrase, os.Stdout, token)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		err = upload(*passphrase, os.Stdin, token)
		if err != nil {
			log.Fatal(err)
		}
		return
	case 2:
		arg0 := flag.Arg(0)
		arg1 := flag.Arg(1)
		if isSHA1Token(arg0) && !isSHA1Token(arg1) {
			fpath := arg1
			token := arg0
			file, err := os.Create(fpath)

			if err != nil {
				log.Fatal(err)
			}
			err = download(*passphrase, file, token)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		if !isSHA1Token(arg0) || isSHA1Token(arg1) {
			fpath := arg0
			token := arg1
			file, err := os.Open(fpath)
			if err != nil {
				log.Fatal(err)
			}
			err = upload(*passphrase, file, token)
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		log.Fatal("No recognized token on command line")
	}
}
