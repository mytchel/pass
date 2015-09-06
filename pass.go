package main

/*
#include <termios.h>
#include <stdio.h>
#include <errno.h>

struct termios tio;

void noecho() {
	struct termios tio;
	tcgetattr(0, &tio);
	tio.c_lflag &= ~ECHO;
	tcsetattr(0, TCSANOW, &tio);
}

void savetermios() {
	tcgetattr(0, &tio);
}

void resettermios() {
	tcsetattr(0, TCSANOW, &tio);
}
*/
import "C"

import (
	"fmt"
	"os"
	"flag"
	"crypto/aes"
	"crypto/cipher"
)

func ReadPassword() []byte {
	fmt.Fprint(os.Stderr, "Enter pass: ")
	
	C.savetermios()
	C.noecho()
	
	data := make([]byte, 32)
	n, err := os.Stdin.Read(data)
	if err != nil {
		panic(err)
	}
	
	C.resettermios()

	i := 0
	for i < n {
		if data[i] == '\n' {
			break
		}
		i++
	}

	if i == n {
		fmt.Errorf("Sorry, password can only be 32 characters.")
		os.Exit(1)
	}
	return data[:i]
}

func clean(a []byte, n int) {
	for i := n; i < len(a); i++ {
		a[i] = 0
	}
}

func Encrypt(block cipher.Block, clear, encrypt *os.File) error {
	in := make([]byte, block.BlockSize())
	out := make([]byte, block.BlockSize())
	
	for {
		n, err := clear.Read(in)
		if n == 0 {
			break
		} else if err != nil {
			return err
		}
		
		clean(in, n)
		
		block.Encrypt(out, in)
		
		_, err = encrypt.Write(out)
		if err != nil {
			return err
		}
	}	
	
	return nil
}

func Decrypt(block cipher.Block, encrypt, clear *os.File) error {
	in := make([]byte, block.BlockSize())
	out := make([]byte, block.BlockSize())
	
	for {
		n, err := encrypt.Read(in)
		if n == 0 {
			break
		} else if err != nil {
			return err
		}
		
		clean(in, n)
		
		block.Decrypt(out, in)
		
		_, err = clear.Write(out)
		if err != nil {
			return err
		}
	}	
	
	return nil
}

func main() {
	var output *os.File
	var err error
	
	decrypt := flag.Bool("d", false, "Decrypt given files to output")
	encrypt := flag.Bool("e", false, "Encrypt given files to output")
	outputPath := flag.String("o", "", "Redirect output to file")
	
	flag.Parse()
	
	if *decrypt && *encrypt || !*decrypt && !*encrypt {
		fmt.Errorf("Please specify one of d or e\n")
		os.Exit(1)
	}
	
	if *outputPath == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(*outputPath)
		if err != nil {
			panic(err)
		}
	}
	
	pass := ReadPassword()
	key := make([]byte, 32)
	copy(key, pass)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	
	for _, arg := range flag.Args() {
		var file *os.File
		
		if arg == "-" {
			file = os.Stdin
		} else {
			file, err = os.Open(arg)
			if err != nil {
				panic(err)
			}
		}
		
		if *encrypt {
			Encrypt(block, file, output)
		} else if *decrypt {
			Decrypt(block, file, output)
		}
	}
}