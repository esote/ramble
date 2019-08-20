package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/majiru/ramble"
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

var server string

func main() {
	server = "http://localhost:8080"

	uuid, err := welcomeHello()

	if err != nil {
		log.Fatal(err)
	}

	if err = welcomeVerify(uuid); err != nil {
		log.Fatal(err)
	}

	uuid, err = sendHello()

	if err != nil {
		log.Fatal(err)
	}

	if err = sendVerify(uuid); err != nil {
		log.Fatal(err)
	}
}
