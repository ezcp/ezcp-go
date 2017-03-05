package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"

	"io"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	ezcpBitcoin = ".ezcp-bitcoin"
	ezcpToken   = ".ezcp-token"
)

func showHelp() {
	fmt.Println("  Premium usage:")
	fmt.Println("    ezcp --bitcoin                   get address for paiement")
	fmt.Println("    ezcp --login <transactionId>     retreive a token and store it")
	fmt.Println("    ezcp <filepath>                  if <filepath> exists, upload the file using the previously stored token")
	fmt.Println("                                     if <filepath> doesn't exist, download the file pointed by previously stored token")

	fmt.Println("\n  Free usage:")
	fmt.Println("    ezcp <filepath> <token>          upload the file using a free token get thanks to the website http://ezcp.io")
	fmt.Println("    ezcp <token> <filepath>          download the file pointed by the token ")

	fmt.Println("\n  More usage:")
	fmt.Println("    cat <file> | ezcp               upload the piped file")
	fmt.Println("    ezcp > file                     download the file to the redirected pipe ")
	fmt.Println(`    ezcp -x "pass phrase" <file>    upload/download the file with encryption`)
}

func urlFromToken(token string, route string) string {
	return "https://api" + string(token[0]) + ".ezcp.io/" + route + "/" + token
}

func isStatusOK(statuscode int) bool {
	return 200 <= statuscode && statuscode < 300
}

func isSHA1Token(token string) bool {
	regex := regexp.MustCompile("^[0-9a-f]{40}$")
	return regex.MatchString(token)
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func readHomeFile(name string) (string, error) {
	path := filepath.Join(homeDir(), name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func getBitcoinAddress() (string, error) {
	address, err := readHomeFile(ezcpBitcoin)
	if err != nil {
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://ezcp.io/bitcoin", nil)
		if err != nil {
			return "", err
		}
		var res *http.Response
		res, err = client.Do(req)
		if err != nil {
			return "", err
		}
		bytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		address = string(bytes)

		err = ioutil.WriteFile(path.Join(homeDir(), ezcpBitcoin), bytes, 0600)
		if err != nil {
			log.Println(err)
		}
	}
	return address, nil
}

func getToken(tx string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://ezcp.io/token/"+tx, nil)
	if err != nil {
		return "", err
	}
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	token := string(bytes)
	err = ioutil.WriteFile(path.Join(homeDir(), ezcpToken), bytes, 0600)
	if err != nil {
		log.Println(err)
	}
	return token, nil
}

func getDurableToken() (string, error) {
	bytes, err := ioutil.ReadFile(path.Join(homeDir(), ezcpToken))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func stdinTerminal() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}
func stdoutTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

func download(pass string, file *os.File, token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlFromToken(token, "download"), nil)
	if err != nil {
		return err
	}

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if !isStatusOK(res.StatusCode) {
		log.Print("ezcp upload status: ", res.StatusCode)
		return err
	}
	var reader io.Reader
	reader, err = crypt(pass, res.Body)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}
	return nil
}
func upload(pass string, file *os.File, token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlFromToken(token, "upload"), nil)
	if err != nil {
		return err
	}
	var reader io.Reader
	reader, err = crypt(pass, file)
	if err != nil {
		panic(err)
	}
	req.Body = ioutil.NopCloser(reader)

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	if !isStatusOK(res.StatusCode) {
		log.Print("ezcp upload status: ", res.StatusCode)
		return err
	}
	return nil
}
