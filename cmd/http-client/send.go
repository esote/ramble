package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/esote/ramble"
)

func sendHello() (uuid string, err error) {
	fmt.Println("Enter public fingerprint:")

	var b bytes.Buffer

	if _, err = b.ReadFrom(os.Stdin); err != nil {
		return
	}

	req := ramble.SendHelloReq{
		Sender: strings.TrimSpace(b.String()),
	}

	b.Reset()

	fmt.Println("Enter conversion UUID (empty for new conversation):")

	if _, err = b.ReadFrom(os.Stdin); err != nil {
		return
	}

	req.Conversation = strings.TrimSpace(b.String())

	b.Reset()

	fmt.Println("Enter recipients' public fingerprints (comma separated)")

	if _, err = b.ReadFrom(os.Stdin); err != nil {
		return
	}

	req.Recipients = strings.Split(strings.TrimSpace(b.String()), ",")

	b.Reset()

	fmt.Println("Enter encrypted message:")

	if _, err = b.ReadFrom(os.Stdin); err != nil {
		return
	}

	req.Message = strings.TrimSpace(b.String())

	data, err := json.Marshal(&req)

	if err != nil {
		return
	}

	body, err := request("/send/hello", data)

	if err != nil {
		return
	}

	var resp ramble.SendHelloResp

	if err = json.Unmarshal(body, &resp); err != nil {
		return
	}

	uuid = resp.UUID

	fmt.Println("Sign nonce with public key:")
	fmt.Println(resp.Nonce)

	return
}

func sendVerify(uuid string) error {
	body, err := verify("/send/verify", uuid)

	if err != nil {
		return err
	}

	var resp ramble.SendVerifyResp

	if err = json.Unmarshal(body, &resp); err != nil {
		return err
	}

	fmt.Printf("Conversation UUID: %s\n", resp.Conversation)

	return nil
}
