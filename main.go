package main

import (
	"flag"
	"fmt"
	"os"
)

var secstorePath *string

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "If no arguments are given then a repl is run.")
}

func initNewSecstore(file *os.File) error {
	var err error
	var pass []byte

	fmt.Fprintln(os.Stderr, "Creating a new secstore...")
	fmt.Fprintln(os.Stderr, "Enter the password to encrypt it with: ")

	if pass, err = GetNewPass(); err != nil {
		return err
	}

	err = EncryptBytes(pass, []byte(nil), file)
	return err
}

func SaveSecstore(store *Secstore) error {
	var err error
	var file *os.File

	if file, err = os.Create(*secstorePath); err != nil {
		return err
	}

	if err = store.EncryptToFile(file); err != nil {
		return err
	}

	file.Close()
	return nil
}

func Exit(store *Secstore) {
	
}

func main() {
	var secstore *Secstore
	var err error
	var file *os.File
	var pass []byte

	secstorePath = flag.String("P", os.Getenv("HOME")+
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
	if pass, err = ReadPassword(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading pass:", err)
		os.Exit(1)
	}

	secstore = new(Secstore)
	secstore.Pass = pass

	if err = secstore.DecryptFile(file); err != nil {
		fmt.Fprintln(os.Stderr, "Error decrypting:", err)
		os.Exit(1)		
	}
	
	file.Close()

	args := flag.Args()
	if len(args) == 0 {
		RunRepl(secstore)
	} else {
		err := EvalCommand(secstore, args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	
	if err := SaveSecstore(secstore); err != nil {
		fmt.Fprintln(os.Stderr, "Error saving:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
