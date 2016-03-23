package main

import (
	"flag"
	"fmt"
	"os"
)

var makePart *string = flag.String("n", "", "Add a new password.")
var makeDir *string = flag.String("m", "", "Add a new directory.")
var show *string = flag.String("s", "", "Show a password.")
var remove *string = flag.String("r", "", "Remove a password.")
var edit *string = flag.String("e", "", "Edit a password.")

var changePass *bool = flag.Bool("P", false, "Change secstore password.")

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "If no arguments are given then it is interpreted as -l '.*'")
}

func getNewPass() []byte {
	var pass1, pass2 []byte
	var good bool = false

	fmt.Fprint(os.Stderr, "Password: ")
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

	return pass1
}

func initNewSecstore(file *os.File) error {
	var err error
	var pass []byte

	fmt.Fprintln(os.Stderr, "Creating a new secstore...")
	fmt.Fprintln(os.Stderr, "Enter the password to encrypt it with: ")

	pass = getNewPass()

	err = EncryptBytes(pass, []byte(nil), file)
	return err
}

func main() {
	var secstore *Secstore
	var err error
	var plain []byte
	var file *os.File

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
	
	secstore, err = ParseSecstore(plain)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)		
	}
	
	if *changePass {
		fmt.Fprintln(os.Stderr, "Changing password...")
		pass = getNewPass()
	} else if len(*makePart) > 0 {
		secstore.MakeNewPart(*makePart)
	} else if len(*makeDir) > 0 {
		secstore.MakeNewDirPart(*makeDir)
	} else if len(*show) > 0 {
		secstore.ShowPart(*show)
	} else if len(*remove) > 0 {
		secstore.RemovePart(*remove)
	} else if len(*edit) > 0 {
		secstore.EditPart(*edit)
	} else {
		secstore.List()
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
