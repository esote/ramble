package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/majiru/ramble"
)

func welcomeHello() (uuid string, err error) {
	fmt.Println("Enter ASCII-armored public key:")

	var b bytes.Buffer

	if _, err = b.ReadFrom(os.Stdin); err != nil {
		return
	}

	req := ramble.WelcomeHelloReq{
		Public: b.String(),
	}

	data, err := json.Marshal(&req)

	if err != nil {
		return
	}

	body, err := request("/welcome/hello", data)

	if err != nil {
		return
	}

	var resp ramble.WelcomeHelloResp

	if err = json.Unmarshal(body, &resp); err != nil {
		return
	}

	uuid = resp.UUID

	fmt.Println("Sign nonce with public key:")
	fmt.Println(resp.Nonce)

	return
}

func welcomeVerify(uuid string) error {
	_, err := verify("/welcome/verify", uuid)

	if err != nil {
		return err
	}

	fmt.Println("The server has welcomed you!")

	return nil
}
