package main

import (
	"fmt"
	"os"
	"flag"
	"crypto/aes"
)

const (
	KeySize = 32
)

var SecstoreStart []byte = []byte("Secstore.\n")

var makeNew *string = flag.String("n", "", "Add a new password.")
var show *string = flag.String("s", "", "Show a password.")
var remove *string = flag.String("r", "", "Remove a password.")
var edit *string = flag.String("e", "", "Edit a password.")
var passwordIn *string = flag.String("P", "/dev/tty", "Where to read the unlock password from.")

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

func decryptFile(pass []byte, file *os.File) ([]byte, error) {
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

func encryptBytes(pass, bytes []byte, file *os.File) error {
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

func initNewSecstore(file *os.File) error {
	var pass1, pass2 []byte
	var good bool = false
	var err error

	fmt.Println("Creating a new secstore...")
	fmt.Print("Enter the password to encrypt it with: ")

	for !good {
		pass1 = ReadPassword()
		fmt.Print("And again: ")
		pass2 = ReadPassword()

		if len(pass1) != len(pass2) {
			good = false
		} else {
			good = true
			for k, _ := range pass1 {
				if pass1[k] != pass2[k] {
					good = false
					break
				}
			}
		}

		if !good {
			fmt.Print("Passwords did not match.\nTry again: ")
		}
	}

	err = encryptBytes(pass1, SecstoreStart, file)
	return err
}

func main() {
	var err error
	var plain []byte
	var file *os.File
	var i int

	secstore := flag.String("p", os.Getenv("HOME") + "/.secstore", "Path to secstore.")

	flag.Parse()

	file, err = os.Open(*secstore)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(*secstore)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else {
				if initNewSecstore(file) != nil {
					fmt.Println(err)
					os.Exit(1)
				} else {
					file.Close()
					os.Exit(0)
				}
			}
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Print("Enter pass: ")
	pass := ReadPassword()

	fmt.Println("pass: ", string(pass))

	plain, err = decryptFile(pass, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i = 0; i < len(SecstoreStart); i++ {
		if plain[i] != SecstoreStart[i] {
			fmt.Println("Could not decrypt secstore")
			os.Exit(1)
		}
	}

	fmt.Println("secstore succesfully decrypted with pass: ", string(pass))
	fmt.Println("secstore:\n")
	fmt.Println(string(plain))

	if len(*makeNew) > 0 {
		fmt.Println("Creating a new password: ", *makeNew)

		newpart := MakeNewPart(*makeNew)
		plain = append(plain, newpart...)

	} else if len(*show) > 0 {
		fmt.Println("showing :", *show)

	} else if len(*remove) > 0 {
		fmt.Println("Removing a password: ", *remove)
	} else if len(*edit) > 0 {
		fmt.Println("Editing a password: ", *edit)
	} else {
		fmt.Println("Showing a list of passwords...")
	}

	file.Close()
	err = os.Remove(*secstore)
	if err != nil {
		panic(err)
	}

	file, err = os.Create(*secstore)
	if err != nil {
		panic(err)
	}

	fmt.Println("encrypt with : ", string(pass))
	err = encryptBytes(pass, plain, file)
	if err != nil {
		panic(err)
	}

	file.Close()
}
