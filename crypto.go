package main

import (
	"crypto/aes"
	"os"
)

const (
	KeySize = 32
)

var SecstoreStart []byte = []byte("Secstore.\n")

func createNewPass(old, bytes []byte) []byte {
	var n []byte
	var i int
	var sum byte

	n = make([]byte, KeySize)

	sum = 0
	for i = 0; i < len(bytes); i++ {
		sum += bytes[i]
	}

	if sum == 0 {
		sum = 1
	}

	for i = 0; i < KeySize; i++ {
		n[i] = old[i] + sum
	}

	return n
}

func DecryptFile(pass []byte, file *os.File) ([]byte, error) {
	var plain, cipher, blockpass []byte
	var plainFull []byte

	plain = make([]byte, aes.BlockSize)
	cipher = make([]byte, aes.BlockSize)
	blockpass = make([]byte, KeySize)

	copy(blockpass, pass)

	for {
		n, err := file.Read(cipher)
		if n != aes.BlockSize {
			break
		} else if err != nil {
			return []byte(nil), err
		}

		conv, err := aes.NewCipher(blockpass)
		if err != nil {
			return []byte(nil), err
		}

		conv.Decrypt(plain, cipher)

		blockpass = createNewPass(blockpass, plain)
		plainFull = append(plainFull, plain...)
	}

	return plainFull, nil
}

func EncryptBytes(pass, bytes []byte, file *os.File) error {
	var plain, cipher, blockpass []byte
	var n, nn int

	plain = make([]byte, aes.BlockSize)
	cipher = make([]byte, aes.BlockSize)
	blockpass = make([]byte, KeySize)

	n = 0

	copy(blockpass, pass)

	for n < len(bytes) {
		nn = copy(plain, bytes[n:])
		for i := nn; i < len(plain); i++ {
			plain[i] = 0
		}

		n += nn

		conv, err := aes.NewCipher(blockpass)
		if err != nil {
			panic(err)
		}

		conv.Encrypt(cipher, plain)

		_, err = file.Write(cipher)
		if err != nil {
			return err
		}

		blockpass = createNewPass(blockpass, plain)
	}

	return nil
}
