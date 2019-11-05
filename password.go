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

	"github.com/peterh/liner"
)

func ReadPassword(line *liner.State) ([]byte, error) {
	s, err := line.PasswordPrompt("Enter pass: ")
	if err != nil {
		return nil, err
	}

	var data []byte = make([]byte, KeySize)
	for k, c := range []byte(s) {
		if c == '\n' {
			fmt.Fprint(os.Stderr, "have new line")
			data[k] = 0
		} else {
			data[k] = c
		}
	}

	return data, nil
}

func GetNewPass(line *liner.State) ([]byte, error) {
	var data []byte = make([]byte, KeySize)
	var pass1, pass2 string
	var good bool = false
	var err error

	if pass1, err = line.PasswordPrompt("New pass: "); err != nil {
		return nil, err
	}

	if pass2, err = line.PasswordPrompt("And again: "); err != nil {
		return nil, err
	}

	if len(pass1) != len(pass2) {
		good = false
	} else {
		good = true
		for k, _ := range pass1 {
			if pass1[k] == '\n' {
				data[k] = 0
				break
			} else if pass1[k] != pass2[k] {
				good = false
				break
			} else {
				data[k] = pass1[k]
			}
		}
	}

	if good {
		return data, nil
	} else {
		return nil, fmt.Errorf("Passwords did not match.")
	}
}
