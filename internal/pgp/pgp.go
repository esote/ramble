package pgp

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"io"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

// TODO (esote): Add internal/pgp function to fetch fingerprint from public key.

// TODO (esote): Add internal/pgp function to verify ASCII armored, encrypted
// messages are really what they seem. Then the check below will not be needed.

const (
	nonceLen = 1024

	// NonceHexLen is the length of a hexadecimal nonce in bytes.
	NonceHexLen = 2 * nonceLen
)

// EncryptArmored encrypts plaintext for recipients by proving a plaintext and
// list of armored public keys. Returns an armored, encrypted PGP message.
//
// The list of recipients is hidden, but the number of recipients can still be
// determined.
func EncryptArmored(public, plain io.Reader) ([]byte, error) {
	key, err := openpgp.ReadArmoredKeyRing(public)

	if err != nil {
		return nil, err
	}

	// Use speculative key IDs to hide encryption recipients. Countermeasure
	// against traffic analysis.
	for _, entity := range key {
		entity.PrimaryKey.KeyId = 0

		for _, subkey := range entity.Subkeys {
			subkey.PublicKey.KeyId = 0
		}
	}

	config := &packet.Config{
		DefaultCipher: packet.CipherAES256,
		DefaultHash:   crypto.SHA512,
	}

	hints := &openpgp.FileHints{
		IsBinary: true,
	}

	var encrypted bytes.Buffer

	wc, err := openpgp.Encrypt(&encrypted, key, nil, hints, config)

	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(wc, plain); err != nil {
		return nil, err
	}

	if err = wc.Close(); err != nil {
		return nil, err
	}

	var armored bytes.Buffer

	wc, err = armor.Encode(&armored, "PGP MESSAGE", nil)

	if err != nil {
		return nil, err
	}

	if _, err = encrypted.WriteTo(wc); err != nil {
		return nil, err
	}

	if err = wc.Close(); err != nil {
		return nil, err
	}

	return armored.Bytes(), nil
}

// NonceHex generates a random nonce encoded as hex.
func NonceHex() (nonce []byte, err error) {
	b := make([]byte, nonceLen)

	if _, err = rand.Read(b); err != nil {
		return
	}

	nonce = make([]byte, NonceHexLen)
	_ = hex.Encode(nonce, b)

	return
}

// VerifyArmoredSig uses an armored public key to verify an armored, detached
// signature.
func VerifyArmoredSig(public, sig, file io.Reader) (bool, error) {
	k, err := openpgp.ReadArmoredKeyRing(public)

	if err != nil {
		return false, err
	}

	_, err = openpgp.CheckArmoredDetachedSignature(k, file, sig)

	return err == nil, err
}
