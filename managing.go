package main

import (
	"fmt"
)

func MakeNewPart(name string) []byte {
	var bytes, part []byte
	
	part = []byte("Hello there.")

	bytes = []byte(name)

	bytes = append(bytes, ':')
	bytes = append(bytes, part...)
	bytes = append(bytes, '\n')

	return bytes
}
