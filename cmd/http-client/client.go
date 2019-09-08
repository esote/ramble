package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/esote/ramble"
)

func request(path string, data []byte) ([]byte, error) {
	uri, err := url.Parse(server + path)

	if err != nil {
		return nil, err
	}

	resp, err := http.Post(uri.String(), "application/json",
		bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

func verify(path, uuid string) ([]byte, error) {
	fmt.Println("Enter nonce detached signature:")

	var b bytes.Buffer

	if _, err := b.ReadFrom(os.Stdin); err != nil {
		return nil, err
	}

	req := ramble.VerifyRequest{
		Signature: b.String(),
		UUID:      uuid,
	}

	data, err := json.Marshal(&req)

	if err != nil {
		return nil, err
	}

	return request(path, data)
}

func shell() error {
	fmt.Println("cmds: d/h/q/s/v/w")

	const help = `d, delete
	Request to delete data from ramble server.
h, help
	Print this help.
q, quit
	Exit this client.
s, send
	Send a message.
v, view
	View messages.
w, welcome
	Introduce yourself to the server.`

	b := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		input, _, err := b.ReadLine()

		if err == io.EOF {
			fmt.Println()
			return nil
		} else if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		input = bytes.TrimSpace(input)

		if len(input) == 0 {
			continue
		}

		var uuid string

		switch string(input) {
		case "d", "delete":
			fmt.Println("not implemented yet")
		case "h", "help":
			fmt.Println(help)
		case "q", "quit":
			return nil
		case "s", "send":
			uuid, err = sendHello()

			if err == nil {
				err = sendVerify(uuid)
			}
		case "v", "view":
			fmt.Println("not implemented yet")
		case "w", "welcome":
			uuid, err = welcomeHello()

			if err == nil {
				err = welcomeVerify(uuid)
			}
		default:
			fmt.Fprintf(os.Stderr, "'%s' is an invalid option\n",
				input)
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

var server string

func main() {
	flag.StringVar(&server, "server", "http://localhost:8080", "server URL")
	flag.Parse()

	if err := shell(); err != nil {
		log.Fatal(err)
	}
}
