package main

import (
	"crypto/aes"
	"fmt"
	"os"
	"strings"
)

var decryptionFuncs = map[string](func([]byte, []byte, *os.File) ([]byte, error)){
	"SecstorePass 0.1": VersionOne,
	"store 02":         VersionTwo,
}

/* Decrypts first block with pass, checks it decrypted alright, then passes
 * along to the proper decryption version. */
func DecryptFile(pass []byte, file *os.File) ([]byte, error) {
	var clear, cipher []byte

	clear = make([]byte, aes.BlockSize)
	cipher = make([]byte, aes.BlockSize)

	n, err := file.Read(cipher)
	if n != aes.BlockSize {
		return []byte(nil), fmt.Errorf("Read error. File too small?")
	} else if err != nil {
		return []byte(nil), err
	}

	conv, err := aes.NewCipher(pass)
	if err != nil {
		return []byte(nil), err
	}

	conv.Decrypt(clear, cipher)

	str := string(clear)
	for header, function := range decryptionFuncs {
		if strings.HasPrefix(str, header) {
			bytes, err := function(pass, clear, file)
			return bytes, err
		}
	}

	return nil, fmt.Errorf("Failed to decrypt.")
}
