package main

import (
	"fmt"
	"regexp"
	"os"
	"os/exec"
	"math/rand"
)

type Part struct {
	Name string
	Data string
	Next *Part
}

type Secstore struct {
	Parts *Part
}

func ParseSecstore(secstoreRaw []byte) *Secstore {
	var i, j, l int
	var line, pname, pdata []byte
	var part, head, prev *Part
	var secstore *Secstore

	head = new(Part)
	head.Next = nil
	prev = head

	l = len(secstoreRaw)
	i = 0
	for i < l {
		for j = i; j < l; j++ {
			if secstoreRaw[j] == 0 {
				break
			}
		}

		if j == l {
			break
		}

		line = secstoreRaw[i:j]

		i = j + 1

		for j = 0; j < len(line); j++ {
			if line[j] == ':' {
				break
			}
		}

		if j == len(line) {
			continue
		}

		pname = line[0:j]
		pdata = line[j+1:]

		part = new(Part)
		part.Name = string(pname)
		part.Data = string(pdata)
		part.Next = nil

		prev.Next = part
		prev = part
	}
	
	secstore = new(Secstore)
	secstore.Parts = head.Next

	return secstore
}

func (store *Secstore) ToBytes() []byte {
	var bytes []byte
	var part *Part

	bytes = make([]byte, len(SecstoreStart))
	copy(bytes, SecstoreStart)

	for part = store.Parts; part != nil; part = part.Next {
		bytes = append(bytes, []byte(part.Name)...)
		bytes = append(bytes, ':')
		bytes = append(bytes, []byte(part.Data)...)
		bytes = append(bytes, 0)
	}

	return bytes
}

func (store *Secstore) FindPart(name string) *Part {
	var part *Part

	regex, err := regexp.Compile("^" + name + "$")
	if err != nil {
		fmt.Println("Error creating regex: ", err)
		return nil
	}

	for part = store.Parts; part != nil; part = part.Next {
		if regex.Match([]byte(part.Name)) {
			return part
		}
	}

	return nil
}

func (store *Secstore) RemovePart(name string) {
	var part, p *Part

	part = store.FindPart(name)
	if part == nil {
		fmt.Println(name, " not found")
		return
	}

	fmt.Println("Removing: ", part.Name)

	if store.Parts == part {
		store.Parts = part.Next
		return
	}

	for p = store.Parts; p != nil; p = p.Next {
		if p.Next == part {
			p.Next = part.Next
		}
	}
}

func (store *Secstore) ShowPart(name string) {
	part := store.FindPart(name)
	if part == nil {
		fmt.Println(name, "not found")
	} else {
		fmt.Println(part.Data)
	}
}

func (store *Secstore) ShowList(pattern string) {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error creating regex: ", err)
		return
	}

	for part := store.Parts; part != nil; part = part.Next {
		if regex.Match([]byte(part.Name)) {
			fmt.Println(part.Name)
		}
	}

}

func (store *Secstore) MakeNewPart(name string) {
	var part *Part
	var err error

	part = store.FindPart(name)
	if part != nil {
		fmt.Println(name, "already exists. Not adding.")
		return
	}

	part = new(Part)

	part.Data, err = OpenEditor("Store your note/password here (remove this).")
	if err != nil {
		fmt.Println("Not adding. Error running editor:", err)
	} else {
		part.Name = name
		fmt.Println("Adding part with name: ", name)
		part.Next = store.Parts
		store.Parts = part
	}
}

func (store *Secstore) EditPart(name string) {
	var part *Part
	var err error
	var data string

	part = store.FindPart(name)
	if part == nil {
		fmt.Println(name, "not found.")
		return
	}

	data, err = OpenEditor(part.Data)

	if err != nil {
		fmt.Println("Not saving. Error running editor:", err)
	} else {
		part.Data = data
	}
}

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
