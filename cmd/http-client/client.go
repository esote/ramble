package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/majiru/ramble"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	var b bytes.Buffer

	for s.Scan() {
		b.WriteString(s.Text() + "\n")
	}

	helloReq := ramble.WelcomeHelloReq{
		Public: b.String(),
	}

	j, err := json.Marshal(&helloReq)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("http://localhost:8080/welcome/hello", "application/json", bytes.NewReader(j))

	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var helloResp ramble.WelcomeHelloResp

	if err = json.Unmarshal(body, &helloResp); err != nil {
		log.Fatal(err)
	}

	fmt.Println(helloResp.Nonce)

	b.Reset()
	s = bufio.NewScanner(os.Stdin)

	for s.Scan() {
		b.WriteString(s.Text() + "\n")
	}

	verifyReq := ramble.WelcomeVerifyReq{
		Signature: b.String(),
		UUID:      helloResp.UUID,
	}

	j, err = json.Marshal(&verifyReq)

	if err != nil {
		log.Fatal(err)
	}

	resp, err = http.Post("http://localhost:8080/welcome/verify", "application/json", bytes.NewReader(j))

	if err != nil {
		log.Fatal(err)
	}

	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
	fmt.Println(string(body))
}
