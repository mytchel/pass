package main

import (
	"fmt"
	"os"
	"flag"
)

var makeNew *string = flag.String("n", "", "Add a new password.")
var show *string = flag.String("s", "", "Show a password.")
var remove *string = flag.String("r", "", "Remove a password.")
var edit *string = flag.String("e", "", "Edit a password.")
var list *string = flag.String("l", "", "List passwords that match a patten.")
var passwordIn *string = flag.String("P", "/dev/tty", "Where to read the unlock password from.")

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "If no arguments are given then it is interpreted as -l '.*'")
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

	err = EncryptBytes(pass1, SecstoreStart, file)
	return err
}

func main() {
	var secstore *Secstore
	var err error
	var plain []byte
	var file *os.File
	var i int

	secstorePath := flag.String("p", os.Getenv("HOME") + "/.secstore", "Path to secstore.")

	flag.Usage = Usage

	flag.Parse()

	file, err = os.Open(*secstorePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(*secstorePath)
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

	plain, err = DecryptFile(pass, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i = 0; i < len(SecstoreStart); i++ {
		if plain[i] != SecstoreStart[i] {
			fmt.Println("Failed to decrypt secstore")
			os.Exit(1)
		}
	}

	secstore = ParseSecstore(plain)

	if len(*makeNew) > 0 {
		secstore.MakeNewPart(*makeNew)
	} else if len(*show) > 0 {
		secstore.ShowPart(*show)
	} else if len(*remove) > 0 {
		secstore.RemovePart(*remove)
	} else if len(*edit) > 0 {
		secstore.EditPart(*edit)
	} else if len(*list) > 0 {
		secstore.ShowList(*list)
	} else {
		secstore.ShowList("")
	}

	file.Close()
	err = os.Remove(*secstorePath)
	if err != nil {
		fmt.Println("Error removing secstore" , *secstorePath, " : ", err)
		os.Exit(1)
	}

	file, err = os.Create(*secstorePath)
	if err != nil {
		fmt.Println("Error recreating secstore" , *secstorePath, " : ", err)
		os.Exit(1)
	}

	plain = secstore.ToBytes()

	err = EncryptBytes(pass, plain, file)
	if err != nil {
		fmt.Println("Error encrypting secstore : ", err)
		os.Exit(1)
	}

	file.Close()
	os.Exit(0)
}
