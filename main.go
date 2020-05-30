package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/peterh/liner"
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "If no arguments are given then a repl is run.")
}

func initNewSecstore(line *liner.State, file *os.File) error {
	var err error
	var pass []byte

	fmt.Fprintln(os.Stderr, "Creating a new secstore...")
	fmt.Fprintln(os.Stderr, "Enter the password to encrypt it with: ")

	if pass, err = GetNewPass(line); err != nil {
		return err
	}

	err = EncryptBytes(pass, []byte(nil), file)
	return err
}

func Exit(store *Secstore) {

}

func main() {
	var secstore *Secstore
	var err error
    var path *string
	var file *os.File
	var pass []byte

	path = flag.String("P", os.Getenv("HOME")+
		"/.secstore", "Path to secstore file.")

	flag.Usage = Usage

	flag.Parse()

	line := liner.NewLiner()
	defer line.Close()

	file, err = os.Open(*path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading secstore")
		if os.IsNotExist(err) {
			file, err = os.Create(*path)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			} else {
				if initNewSecstore(line, file) != nil {
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

	file.Close()

	if pass, err = ReadPassword(line); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading pass:", err)
		os.Exit(1)
	}

	secstore = new(Secstore)
	secstore.Pass = pass
	secstore.Path = path

	if err = secstore.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "Error decrypting:", err)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) == 0 {
		RunRepl(secstore, line)
	} else {
		err := EvalCommand(secstore, line, args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	os.Exit(0)
}
