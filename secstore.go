package main

import (
	"fmt"
	"os"
)

type Secstore struct {
	Pass     []byte
	Path     *string
	
	Saved    bool
	rootPart *Part
	Pwd      *Part
}

func (store *Secstore) Load() error {
	var file *os.File
	var err error
	var plain []byte

	file, err = os.Open(*store.Path)
	if err != nil {
		return err
	}

	defer file.Close()

	plain, err = DecryptFile(store.Pass, file)
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
	store.Saved = true

	return nil
}

func (store *Secstore) Save() error {
	var bytes []byte = []byte(nil)
    var err error
	var file *os.File

	if file, err = os.Create(*store.Path); err != nil {
        return err
	}

	for part := store.rootPart.SubParts; part != nil; part = part.Next {
		bytes = append(bytes, part.ToBytes()...)
	}

    err = EncryptBytes(store.Pass, bytes, file)
	
	file.Close()

	if err == nil {
        store.Saved = true
    }

    return err
}

func (store *Secstore) addPart(path []string) (*Part, error) {
	var part, parent *Part

	if part = store.Pwd.FindSub(path); part != nil {
		return nil, fmt.Errorf("'%s' already exists.", path)
	}

	ppath := path[0 : len(path)-1]
	name := path[len(path)-1]

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
