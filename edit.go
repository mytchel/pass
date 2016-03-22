package main

import (
	"fmt"
	"os"
	"os/exec"
	"math/rand"
)

func findValidTempFile(prefix string) string {
	var err error

	path := []byte(os.TempDir())
	path = append(path, byte('/'))
	path = append(path, []byte(prefix)...)

	rand.Seed(int64(rand.Int()))
	for {
		r := rand.Int() % 25 + int('a')

		path = append(path, byte(r))

		_, err = os.Stat(string(path))
		if err != nil && os.IsNotExist(err) {
			break
		}
	}

	return string(path)
}

func OpenEditor(data string) (string, error) {
	var tmpPath string
	var file *os.File
	var err error
	var bytes, newData []byte
	var cmd *exec.Cmd

	tmpPath = findValidTempFile(".pass-")

	file, err = os.Create(tmpPath)
	if err != nil {
		fmt.Println("err: ", err)
		return "", err
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		return "", err
	}

	file.Close()

	/* Open editor */
	
	editor := os.Getenv("EDITOR")
	if len(editor) == 0 {
		editor = "vi"
	}

	cmd = exec.Command(editor, tmpPath)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("got error: ", err)
		return "", err
	}

	/* Read in new data */

	file, err = os.Open(tmpPath)
	if err != nil {
		fmt.Println("Error reopening tmp file: ", tmpPath)
		return "", err
	}

	bytes = make([]byte, 16)

	for {
		n, err := file.Read(bytes)
		if n == 0 {
			break
		} else if err != nil {
			return "", err
		}

		newData = append(newData, bytes[:n]...)
	}

	err = os.Remove(tmpPath)
	if err != nil {
		fmt.Println("Error removing temp file!! ", err)
		fmt.Println("Please remove: ", tmpPath)
	}

	return string(newData), nil
}
