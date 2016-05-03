package main

/*
#include <termios.h>
#include <stdio.h>
#include <errno.h>

struct termios tio_save;

void setnoecho(int fd) {
	struct termios tio;
	tcgetattr(fd, &tio);
	tio.c_lflag &= ~ECHO;
	tcsetattr(fd, TCSANOW, &tio);
}

void savetermios(int fd) {
	tcgetattr(fd, &tio_save);
}

void resettermios(int fd) {
	tcsetattr(fd, TCSANOW, &tio_save);
}
*/
import "C"

import (
	"fmt"
	"os"
)

func ReadPassword() ([]byte, error) {
	var err error
	var data []byte = make([]byte, KeySize)
	var tty *os.File

	if tty, err = os.Open("/dev/tty"); err != nil {
		return nil, err
	}

	C.savetermios(C.int(os.Stdin.Fd()))
	C.setnoecho(C.int(os.Stdin.Fd()))

	_, err = tty.Read(data)
	
	C.resettermios(C.int(os.Stdin.Fd()))
	
	if err != nil {
		return nil, err
	}


	tty.Close()

	fmt.Fprintln(os.Stderr)

	for k, c := range data {
		if c == '\n' {
			data[k] = 0
		}
	}

	return data, nil
}

func GetNewPass() ([]byte, error) {
	var pass1, pass2 []byte
	var good bool = false
	var err error

	fmt.Fprint(os.Stderr, "New Password: ")
	for !good {
		if pass1, err = ReadPassword(); err != nil {
			return nil, err
		}

		fmt.Fprint(os.Stderr, "And again: ")
		if pass2, err = ReadPassword(); err != nil {
			return nil, err
		}

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

	return pass1, nil
}

