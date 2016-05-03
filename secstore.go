package main

import (
	"fmt"
	"os"
	"strings"
	"crypto/rand"
)

type Secstore struct {
	rootPart *Part
	Pwd *Part
	Pass []byte
}

func (store *Secstore) DecryptFile(file *os.File) error {
	var err error
	var plain []byte

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

	return nil
}

func (store *Secstore) EncryptToFile(file *os.File) error {
	var bytes []byte = []byte(nil)

	for part := store.rootPart.SubParts; part != nil; part = part.Next {
		bytes = append(bytes, part.ToBytes()...)
	}

	return EncryptBytes(store.Pass, bytes, file)
}

func ChangePass(store *Secstore, args []string) error {
	var pass []byte
	var err error

	if pass, err = GetNewPass(); err != nil {
		return err
	}

	store.Pass = pass
	return nil
}

func ChangeDir(store *Secstore, args []string) error {
	var n *Part

	if len (args) == 0 {
		store.Pwd = store.rootPart
		return nil
	} else if len(args) > 1 {
		return fmt.Errorf("usage: cd [dir]")
	} else {
		path := strings.Split(args[0], "/")
		if n = store.Pwd.FindSub(path); n != nil {
			store.Pwd = n
			return nil
		} else {
			return fmt.Errorf("'%s' does not exist.", args[0])
		}
	}
}

func RemovePart(store *Secstore, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: rm part1 part2...")
	} else {
		for _, arg := range args {
			path := strings.Split(arg, "/")
			part := store.Pwd.FindSub(path)

			if part == nil {
				return fmt.Errorf("'%s' does not exist.", args[0])
			} else if part.Parent == nil {
				return fmt.Errorf("Not removing root dir.")
			}
			
			if err := part.Parent.RemovePart(part); err != nil {
				return err
			}
		}
		return nil
	}
}

func ShowPart(store *Secstore, args []string) error {
	var path []string
	
	if len(args) > 0 {
		path = strings.Split(args[0], "/")
	} else {
		path = []string(nil)
	}

	part := store.Pwd.FindSub(path)
	if part == nil {
		return fmt.Errorf("'%s' does not exist.", args[0])
	} else {
		part.Print()
		return nil
	}
}

func EditPart(store *Secstore, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: edit part")
	}

	path := strings.Split(args[0], "/")

	part := store.Pwd.FindSub(path)
	if part == nil {
		return fmt.Errorf("'%s' does not exist.", args[0])
	} else if part.Type == TypeDir {
		return fmt.Errorf("'%s' is a directory.", args[0])
	} else {
		if data, err := OpenEditor(part.Data); err != nil {
			return fmt.Errorf("Not saving. Error running editor:", err)
		} else {
			part.Data = data
			return nil
		}
	}
}

func AddDataPart(store *Secstore, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: add name")
	}

	path := strings.Split(args[0], "/")
	if part, err := store.addPart(path); err == nil {
		part.Type = TypeData
		part.Data, _ = OpenEditor(randomPass())
		return nil
	} else {
		return err
	}
}

func AddDirPart(store *Secstore, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mkdir name")
	}

	path := strings.Split(args[0], "/")
	if part, err := store.addPart(path); err != nil {
		part.Type = TypeDir
		part.SubParts = nil
		return nil
	} else {
		return err
	}
}

func MovePart(store *Secstore, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: mv part1 dest")
	}

	return nil	
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


func randomPass() string {
	var sum, r int
	var b []byte
	var err error

	b = make([]byte, 24)
	_, err = rand.Read(b)

	if err != nil {
		return "Error generating random bytes!"
	} 

	sum = 0
	for i := 0; i < len(b); i++ {
		sum += int(b[i])
		r = sum % 3
		switch (r) {
		case 0:
			b[i] = 'a' + b[i] % 26
		case 1:
			b[i] = 'A' + b[i] % 26
		case 2:
			b[i] = '0' + b[i] % 10
		}
	}

	return string(b)
}

