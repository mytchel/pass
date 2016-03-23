package main

import (
	"fmt"
	"os"
)

type Secstore struct {
	partRoot *Part
}

func ParseSecstore(bytes []byte) (*Secstore, error) {
	var err error
	var secstore *Secstore

	secstore = new(Secstore)

	secstore.partRoot = new(Part)
	secstore.partRoot.Name = "/"
	
	secstore.partRoot.SubParts, _, err = ParseParts(bytes)
	if err != nil {
		return nil, err
	}
	
	return secstore, nil
}

func (store *Secstore) ToBytes() []byte {
	var bytes []byte = []byte(nil)

	for part := store.partRoot.SubParts; part != nil; part = part.Next {
		bytes = append(bytes, part.ToBytes()...)
	}

	return bytes
}

func (store *Secstore) FindPart(name string) *Part {
	return store.partRoot.FindSub(name)
}

func (store *Secstore) RemovePart(name string) {
	var part, p *Part

	part = store.FindPart(name)
	if part == nil {
		fmt.Println(name, " not found")
		return
	}

	fmt.Println("Removing: ", part.Name)

	for p = store.partRoot; p != nil; p = p.Next {
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
		part.Print()
	}
}

func (store *Secstore) List() {
	store.partRoot.Print()
}

func (store *Secstore) EditPart(name string) {
	var part *Part
	var err error
	var data string

	part = store.FindPart(name)
	if part == nil {
		fmt.Println(name, "not found.")
	} else if part.Data == "" {
		fmt.Println(name, "is a directory.")
	} else {
		data, err = OpenEditor(part.Data)
		if err != nil {
			fmt.Println("Not saving. Error running editor:", err)
		} else {
			part.Data = data
		}
	}
}

func splitLast(s string, sep rune) (main, last string) {
	for i := len(s) - 1; i >= 0; i-- {
		if rune(s[i]) == sep {
			return s[:i], s[i+1:]
		}
	}

	return "", s
}

func (store *Secstore) addPart(fpath string) (*Part, error) {
	var part, parent *Part
	var path, name string

	part = store.FindPart(fpath)
	if part != nil {
		return nil, fmt.Errorf("%s already exists", fpath)
	}
	
	path, name = splitLast(fpath, '/')
	if len(path) > 0 {
		parent = store.partRoot.FindSub(path)
	} else {
		parent = store.partRoot
	}

	part = new(Part)
	part.Name = name

	part.Next = parent.SubParts
	parent.SubParts = part

	return part, nil
}

func (store *Secstore) MakeNewPart(name string) {
	part, err := store.addPart(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	part.Data, _ = OpenEditor("Store your note/password here (remove this).")
	fmt.Println("Adding password:", name)
}

func (store *Secstore) MakeNewDirPart(name string) {
	part, err := store.addPart(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	part.Data = ""
	part.SubParts = nil

	fmt.Println("Adding directory:", name)
}
