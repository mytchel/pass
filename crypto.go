package main

import (
	"fmt"
	"os"
	"crypto/aes"
)

const (
	KeySize = 32
)

var SecstoreStart []byte = []byte("SecstorePass 0.1")

func createNewPass(oldKey, bytes []byte) []byte {
	var newKey, both []byte
	var i, sum, start int

	both = make([]byte, len(oldKey) + len(bytes))
	copy(both, oldKey)
	copy(both[len(oldKey):], bytes)

	sum = 0
	for i = 0; i < len(both); i++ {
		sum += int(both[i])
	}

	start = sum % (len(both) - KeySize)

	newKey = both[start:(start + KeySize)]

	for i = 0; i < KeySize; i++ {
		newKey[i] = byte(int(newKey[i]) + sum)
	}

	return newKey
}

func DecryptFile(pass []byte, file *os.File) ([]byte, error) {
	var clear, cipher, blockpass []byte
	var plain []byte

	clear = make([]byte, aes.BlockSize)
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

		conv.Decrypt(clear, cipher)

		blockpass = createNewPass(blockpass, clear)
		plain = append(plain, clear...)
	}

	for i := 0; i < aes.BlockSize; i++ {
		if SecstoreStart[i] != plain[i] {
			return nil, fmt.Errorf("Failed to decrypt")
		}
	}

	return plain[aes.BlockSize:], nil
}

func EncryptBytes(pass, bytes []byte, file *os.File) error {
	var plain, cipher, blockpass []byte
	var n, nn int

	plain = make([]byte, aes.BlockSize)
	cipher = make([]byte, aes.BlockSize)
	blockpass = make([]byte, KeySize)

	n = 0

	copy(blockpass, pass)
	copy(plain, SecstoreStart)

	for {
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

		blockpass = createNewPass(blockpass, plain)

		nn = copy(plain, bytes[n:])
		for i := nn; i < len(plain); i++ {
			plain[i] = 0
		}

		n += nn
	}

	return nil
}
