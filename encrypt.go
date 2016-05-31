package main

import (
	"os"
	"crypto/aes"
	"crypto/rand"
)

var secstoreStart []byte = []byte("store 02")

func EncryptBytes(pass, bytes []byte, file *os.File) error {
	var plain, cipher, blockpass []byte
	var n, nn int

	plain = make([]byte, 16)
	cipher = make([]byte, 16)
	blockpass = make([]byte, KeySize)

	n = 0

	copy(blockpass, pass)
	copy(plain, secstoreStart)

	for {
		if _, err := rand.Read(plain[8:]); err != nil {
			panic(err)
		}

		conv, err := aes.NewCipher(blockpass)
		if err != nil {
			panic(err)
		}

		conv.Encrypt(cipher, plain)

		_, err = file.Write(cipher)
		if err != nil {
			return err
		}

		if n >= len(bytes) {
			break
		}

		blockpass = TwoCreateNewPass(blockpass, plain)

		nn = copy(plain[:8], bytes[n:])
		for i := nn; i < 8; i++ {
			plain[i] = 0
		}

		n += 8
	}

	return nil
}
