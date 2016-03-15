package main

/*
#include <termios.h>
#include <stdio.h>
#include <errno.h>

struct termios tio;

void noecho(int fd) {
	struct termios tio;
	tcgetattr(fd, &tio);
	tio.c_lflag &= ~ECHO;
	tcsetattr(fd, TCSANOW, &tio);
}

void savetermios(int fd) {
	tcgetattr(fd, &tio);
}

void resettermios(int fd) {
	tcsetattr(fd, TCSANOW, &tio);
}
*/
import "C"

import (
	"fmt"
	"os"
)

func ReadPassword() []byte {
	var err error
	
	tty, err := os.Open(*passwordIn)
	if err != nil {
		panic(err)
	}
	
	C.savetermios(C.int(tty.Fd()))
	C.noecho(C.int(tty.Fd()))
	
	data := make([]byte, KeySize)
	_, err = tty.Read(data)
	tty.Close()
	if err != nil {
		panic(err)
	}
	
	C.resettermios(C.int(tty.Fd()))
	
	fmt.Print("\n")

	for k, c := range data {
		if c == '\n' {
			data[k] = 0
		}
	}

	return data
}

