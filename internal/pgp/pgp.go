package pgp

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"time"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
)

const (
	encType           = "PGP MESSAGE"
	fingerprintHexLen = 2 * 20
	nonceLen          = 1024

	// NonceHexLen is the length of a hexadecimal nonce in bytes.
	NonceHexLen = 2 * nonceLen
)

// EncryptArmored encrypts plaintext for one recipient by proving a plaintext
// and an armored public key. Returns an armored, encrypted PGP message.
//
// The recipient is hidden using a speculative key ID.
func EncryptArmored(public, plain io.Reader) ([]byte, error) {
	key, err := openpgp.ReadArmoredKeyRing(public)

	if err != nil {
		return nil, err
	}

	// Use speculative key IDs to countermeasure traffic analysis.
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

	wc, err = armor.Encode(&armored, encType, nil)

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

// FingerprintArmored gets the fingerprint from an armored public key.
func FingerprintArmored(public io.Reader) ([]byte, error) {
	key, err := openpgp.ReadArmoredKeyRing(public)

	if err != nil {
		return nil, err
	}

	return key[0].PrimaryKey.Fingerprint[:], nil
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
// signature. Returns the signature creation time.
func VerifyArmoredSig(public, sig, file io.Reader) (t time.Time, err error) {
	k, err := openpgp.ReadArmoredKeyRing(public)

	if err != nil {
		return
	}

	// Read into buffer since we need to read it twice.
	buffer, err := ioutil.ReadAll(sig)

	if err != nil {
		return
	}

	reader := bytes.NewReader(buffer)

	// Takes care of the actual signature validation.
	_, err = openpgp.CheckArmoredDetachedSignature(k, file, reader)

	if err != nil {
		return
	}

	reader = bytes.NewReader(buffer)
	blk, err := armor.Decode(reader)

	if err != nil {
		return
	}

	p, err := packet.NewReader(blk.Body).Next()

	if err != nil {
		return
	}

	// Return the signature creation time.
	switch sig := p.(type) {
	case *packet.Signature:
		return sig.CreationTime, nil
	case *packet.SignatureV3:
		return sig.CreationTime, nil
	default:
		return t, errors.New("incorrect packet type")
	}
}

// VerifyEncryptedArmored tries to validate that input is indeed an armored,
// encrypted PGP message.
func VerifyEncryptedArmored(input io.Reader) (bool, error) {
	blk, err := armor.Decode(input)

	if err != nil {
		return false, err
	}

	if blk.Type != encType {
		return false, errors.New("incorrect block type")
	}

	p, err := packet.NewReader(blk.Body).Next()

	if err != nil {
		return false, err
	}

	switch p.(type) {
	case *packet.EncryptedKey:
		return true, nil
	default:
		return false, errors.New("incorrect packet type")
	}
}

var reHex = regexp.MustCompile("^[a-fA-F0-9]+$")

// VerifyHexFingerprint does rough checks to see if the input fingerprint is
// valid.
func VerifyHexFingerprint(fingerprint string) bool {
	return len(fingerprint) == fingerprintHexLen &&
		reHex.MatchString(fingerprint)
}

// VerifyPublicArmored tries to validate that input is indeed an armored public
// key.
func VerifyPublicArmored(input io.Reader) (bool, error) {
	blk, err := armor.Decode(input)

	if err != nil {
		return false, err
	}

	if blk.Type != openpgp.PublicKeyType {
		return false, errors.New("incorrect block type")
	}

	p, err := packet.NewReader(blk.Body).Next()

	if err != nil {
		return false, err
	}

	switch p.(type) {
	case *packet.PublicKey:
		return true, nil
	default:
		return false, errors.New("incorrect packet type")
	}
}
