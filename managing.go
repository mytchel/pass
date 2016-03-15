package main

import (
	"fmt"
	"regexp"
)

type Part struct {
	Name string
	Data string
}

type Secstore struct {
	Parts []*Part
}

func ParseSecstore(secstoreRaw []byte) *Secstore {
	var i, j, l int
	var line, pname, pdata []byte
	var part *Part
	var secstore *Secstore

	secstore = new(Secstore)

	l = len(secstoreRaw)
	i = 0
	for i < l {
		for j = i; j < l; j++ {
			if secstoreRaw[j] == '\n' {
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

		secstore.Parts = append(secstore.Parts, part)		
	}

	return secstore
}

func (store *Secstore) ToBytes() []byte {
	var bytes []byte

	bytes = make([]byte, len(SecstoreStart))
	copy(bytes, SecstoreStart)

	for _, part := range store.Parts {
		bytes = append(bytes, []byte(part.Name)...)
		bytes = append(bytes, ':')
		bytes = append(bytes, []byte(part.Data)...)
		bytes = append(bytes, '\n')
	}

	return bytes
}

func (store *Secstore) FindPart(name string) (*Part, error) {
	var part *Part

	regex, err := regexp.Compile(name)
	if err != nil {
		return nil, err
	}

	for _, part = range(store.Parts) {
		if regex.Match([]byte(part.Name)) {
			return part, nil
		}
	}

	return nil, nil
}

func (store *Secstore) MakeNewPart(name string) {
	var part *Part

	part, _ = store.FindPart(name)
	if part != nil {
		fmt.Println(name, "already exists. Not adding.")
		return
	}

	part = new(Part)

	part.Name = name
	part.Data = "Hello there."

	fmt.Println("adding part with name: ", name)

	store.Parts = append(store.Parts, part)
}

func (store *Secstore) ShowPart(name string) {
	part, err := store.FindPart(name)
	if err != nil {
		panic(err)
	} else if part == nil {
		fmt.Println(name, "not found")
	} else {
		fmt.Println(part.Data)
	}
}

func (store *Secstore) RemovePart(name string) {
	fmt.Println("Removing: ", name)

}

func (store *Secstore) EditPart(name string) {

}

func (store *Secstore) ShowList(pattern string) {
	fmt.Println("Show list of passwords matching: ", pattern)

	regex, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error creating regex: ", err)
		return
	}

	for _, part := range(store.Parts) {
		if regex.Match([]byte(part.Name)) {
			fmt.Println(part.Name)
		}
	}

}
