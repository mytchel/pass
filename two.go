package main

import (
	"os"
	"crypto/aes"
)

func TwoCreateNewPass(oldKey, bytes []byte) []byte {
	var newKey []byte
	var i int

	newKey = make([]byte, KeySize)
	copy(newKey, oldKey)

	for i = 0; i < KeySize; i++ {
		newKey[i] += bytes[i % len(bytes)]		
	}

	return newKey
}

func VersionTwo(pass, clear []byte, file *os.File) ([]byte, error) {
	var plain, cipher, blockpass []byte

	cipher = make([]byte, aes.BlockSize)
	blockpass = make([]byte, KeySize)

	copy(blockpass, pass)

	for {
		blockpass = TwoCreateNewPass(blockpass, clear)

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

		plain = append(plain, clear[:8]...)
	}

	return plain, nil
}

