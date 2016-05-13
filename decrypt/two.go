package decrypt

import (
	"os"
	"crypto/aes"
)

func twoCreateNewPass(oldKey, bytes []byte) []byte {
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

func VersionTwo(pass, clear []byte, file *os.File) ([]byte, error) {
	var plain, cipher, blockpass []byte

	cipher = make([]byte, aes.BlockSize)
	blockpass = make([]byte, KeySize)

	copy(blockpass, pass)

	for {
		blockpass = twoCreateNewPass(blockpass, clear)

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

		plain = append(plain, clear...)
	}

	return plain, nil
}

