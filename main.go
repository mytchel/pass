package main

import (
	"flag"
	"fmt"
	"os"
)

var makeNew *string = flag.String("n", "", "Add a new password.")
var show *string = flag.String("s", "", "Show a password.")
var remove *string = flag.String("r", "", "Remove a password.")
var edit *string = flag.String("e", "", "Edit a password.")
var list *string = flag.String("l", "", "List passwords that match a patten.")

var dump *bool = flag.Bool("D", false, "Dump secstore to stdout unencrypted.")
var read *string = flag.String("R", "", "Read a clear text secstore and add it to the secstore.")

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

	fmt.Fprintln(os.Stderr, "Creating a new secstore...")
	fmt.Fprint(os.Stderr, "Enter the password to encrypt it with: ")

	for !good {
		pass1 = ReadPassword()
		fmt.Fprint(os.Stderr, "And again: ")
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
			fmt.Fprint(os.Stderr, "Passwords did not match.\nTry again: ")
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
	var start int

	secstorePath := flag.String("p", os.Getenv("HOME")+
		"/.secstore", "Path to secstore file.")

	flag.Usage = Usage

	flag.Parse()

	file, err = os.Open(*secstorePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading secstore")
		if os.IsNotExist(err) {
			file, err = os.Create(*secstorePath)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				if initNewSecstore(file) != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				} else {
					file.Close()
					os.Exit(0)
				}
			}
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	fmt.Fprint(os.Stderr, "Enter pass: ")
	pass := ReadPassword()

	plain, err = DecryptFile(pass, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	
	file.Close()

	for start = 0; start < len(SecstoreStart); start++ {
		if plain[start] != SecstoreStart[start] {
			fmt.Fprintln(os.Stderr, "Failed to decrypt secstore")
			os.Exit(1)
		}
	}

	if *dump {
		fmt.Println(string(plain[start:]))
		os.Exit(0)
	} else if len(*read) > 0 {
		adding, err := os.Open(*read)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file to add: ", err)
			os.Exit(1)
		}

		bytes := make([]byte, 16)
		for {
			n, err := adding.Read(bytes)
			if err != nil {
				break
			}

			plain = append(plain, bytes[:n]...)
		}

		adding.Close()
	}

	secstore = ParseSecstore(plain[start:])
	
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
		secstore.ShowList(".*")
	}

	file, err = os.Create(*secstorePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error recreating secstore", *secstorePath, " : ", err)
		os.Exit(1)
	}

	plain = secstore.ToBytes()

	err = EncryptBytes(pass, plain, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encrypting secstore : ", err)
		os.Exit(1)
	}

	file.Close()
	os.Exit(0)
}
