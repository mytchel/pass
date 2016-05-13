package main

import (
	"fmt"
	"os"

	"lackname.org/pass/decrypt"
)

type Secstore struct {
	rootPart *Part
	Pwd *Part
	Pass []byte
}

func (store *Secstore) DecryptFile(file *os.File) error {
	var err error
	var plain []byte

	plain, err = decrypt.DecryptFile(store.Pass, file)
	if err != nil {
		return err
	}

	store.rootPart = new(Part)
	store.rootPart.Type = TypeDir
	store.rootPart.Name = "/"
	
	store.rootPart.SubParts, _, err = ParseParts(plain, store.rootPart)
	if err != nil {
		return err
	}
	
	store.Pwd = store.rootPart

	return nil
}

func (store *Secstore) EncryptToFile(file *os.File) error {
	var bytes []byte = []byte(nil)

	for part := store.rootPart.SubParts; part != nil; part = part.Next {
		bytes = append(bytes, part.ToBytes()...)
	}

	return EncryptBytes(store.Pass, bytes, file)
}

func (store *Secstore) addPart(path []string) (*Part, error) {
	var part, parent *Part

	if part = store.Pwd.FindSub(path); part != nil {
		return nil, fmt.Errorf("'%s' already exists.", path)
	}
	
	ppath := path[0:len(path) - 1]
	name := path[len(path) - 1]

	if parent = store.Pwd.FindSub(ppath); parent == nil {
		return nil, fmt.Errorf("Parent '%s' does not exist.", ppath)
	}

	part = new(Part)
	part.Name = name
	part.Parent = parent

	part.Next = parent.SubParts
	parent.SubParts = part

	return part, nil
}

