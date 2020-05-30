package main

import (
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/peterh/liner"
)

/* Be careful to not start any commands with q. */
var Commands = map[string](func(*Secstore, *liner.State, []string) error){
	"chpass": ChangePass,
	"add":    AddDataPart,
	"mkdir":  AddDirPart,
	"show":   ShowPart,
	"ls":     ShowPart,
	"rm":     RemovePart,
	"edit":   EditPart,
	"cd":     ChangeDir,
	"mv":     MovePart,
	"help":   Help,
	"save":   Save,
	"quit":   Quit,
}

func randomPass() string {
	var sum, r int
	var b []byte
	var err error

	b = make([]byte, 25)
	_, err = rand.Read(b)

	if err != nil {
		return "Error generating random bytes!"
	}

	sum = 0
	for i := 0; i < len(b) - 1; i++ {
		sum += int(b[i])
		r = sum % 3
		switch r {
		case 0:
			b[i] = 'a' + b[i]%26
		case 1:
			b[i] = 'A' + b[i]%26
		case 2:
			b[i] = '0' + b[i]%10
		}
	}

	b[24] = '\n'

	return string(b)
}

func ChangePass(store *Secstore, line *liner.State, args []string) error {
	var pass []byte
	var err error

	if pass, err = GetNewPass(line); err != nil {
		return err
	}

	store.Pass = pass
	store.Saved = false

	return nil
}

func ChangeDir(store *Secstore, line *liner.State, args []string) error {
	var n *Part

	if len(args) == 0 {
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

func RemovePart(store *Secstore, line *liner.State, args []string) error {
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
        
            store.Saved = false
		}

		return nil
	}
}

func ShowPart(store *Secstore, line *liner.State, args []string) error {
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

func EditPart(store *Secstore, line *liner.State, args []string) error {
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
			return fmt.Errorf("Not saving. Error running editor: %s", err)
		} else {
			part.Data = data
            store.Saved = false
			return nil
		}
	}
}

func AddDataPart(store *Secstore, line *liner.State, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: add name")
	}

	path := strings.Split(args[0], "/")
	if part, err := store.addPart(path); err == nil {
		part.Type = TypeData
		part.Data, _ = OpenEditor(randomPass())
        store.Saved = false
		return nil
	} else {
		return err
	}
}

func AddDirPart(store *Secstore, line *liner.State, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mkdir name")
	}

	path := strings.Split(args[0], "/")
	if part, err := store.addPart(path); err == nil {
		part.Type = TypeDir
		part.SubParts = nil
        store.Saved = false
		return nil
	} else {
		return err
	}
}

func MovePart(store *Secstore, line *liner.State, args []string) error {
	var old, dest *Part
	var err error

	if len(args) != 2 {
		return fmt.Errorf("usage: mv old new")
	}

	oldPath := strings.Split(args[0], "/")
	if old = store.Pwd.FindSub(oldPath); old == nil {
		return fmt.Errorf("'%s' not found.", args[0])
	}

	destPath := strings.Split(args[1], "/")

	dest = store.Pwd.FindSub(destPath)
	if dest != nil && dest.Type == TypeDir {
		destPath = append(destPath, old.Name)
		dest = nil
	}

	if dest == nil {
		if dest, err = store.addPart(destPath); err != nil {
			return err
		}
	}

	old.Parent.RemovePart(old)
	dest.Type = old.Type
	dest.Data = old.Data
	dest.SubParts = old.SubParts

    store.Saved = false

	return nil
}

func Help(store *Secstore, line *liner.State, args []string) error {
	fmt.Print("Commands can be shortened, eg: quit can be called with q.\n\n")
	fmt.Println("chpass\t\tChange encryption password.")
	fmt.Println("add name\tAdd a new password and start editing it.")
	fmt.Println("mkdir name\tMake a new directory.")
	fmt.Println("show name\tShow password or directory.")
	fmt.Println("ls name\t\tAlias for show.")
	fmt.Println("rm name\t\tRemove a password or directory.")
	fmt.Println("edit name\tEdit a password.")
	fmt.Println("cd dir\t\tChange directory.")
	fmt.Println("mv old new\tMove a password or directory from new to old.")
	fmt.Println("save\t\tSave.")
	fmt.Println("help\t\tPrint this.")
	fmt.Println("quit\t\tSave and exit.")
	return nil
}

func Save(store *Secstore, line *liner.State, args []string) error {
	return store.Save()
}

func Quit(store *Secstore, line *liner.State, args []string) error {
	return fmt.Errorf("Quit should never actually be called!")
}

func MatchCommand(cmd string) (func(*Secstore, *liner.State, []string) error, error) {
	var f func(*Secstore, *liner.State, []string) error
	var found bool = false

	for c, g := range Commands {
		if strings.HasPrefix(c, cmd) {
			if found {
				return nil, fmt.Errorf("'%s' matched mulitple commands.", cmd)
			} else {
				found = true
				f = g
			}
		}
	}

	if f == nil {
		return nil, fmt.Errorf("'%s' matched no commands.", cmd)
	} else {
		return f, nil
	}
}
