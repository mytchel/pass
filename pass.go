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

func ReadPassword() []byte {
	var err error
	var data []byte = make([]byte, KeySize)

	C.savetermios(C.int(os.Stdin.Fd()))
	C.setnoecho(C.int(os.Stdin.Fd()))

	_, err = os.Stdin.Read(data)
	if err != nil {
		panic(err)
	}

	C.resettermios(C.int(os.Stdin.Fd()))

	fmt.Fprintln(os.Stderr)

	for k, c := range data {
		if c == '\n' {
			data[k] = 0
		}
	}

	return data
}
