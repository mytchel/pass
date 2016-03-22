package main

import (
	"fmt"
)

type Part struct {
	Name string

	Data string
	SubParts *Part

	Next *Part
}

type Secstore struct {
	Parts *Part
}

func parsePart(bytes []byte) (*Part, int, error) {
	var j, k int
	var part *Part
	var err error

	part = new(Part)

	for j = 0; j < len(bytes) && bytes[j] != 0; j++ {}
	if j == len(bytes) {
		return nil, j, fmt.Errorf("Error parsing secstore: Reached end.")
	}

	part.Name = string(bytes[:j])
	
	for k = j + 1; k < len(bytes) && bytes[k] != 0; k++ {}
	if k == len(bytes) {
		return nil, k, fmt.Errorf("Error parsing secstore: Reached end.")
	}
	
	/* Data part */
	if k > j + 1 {
		part.Data = string(bytes[j+1:k])

	/* Sub tree */
	} else {
		part.Data = ""
		part.SubParts, k, err = parseParts(bytes[k:])
		if err != nil {
			return nil, k, err
		}
	}

	return part, k, nil
}

func parseParts(bytes []byte) (*Part, int, error) {
	var i, j int
	var part, head, prev *Part
	var err error

	head = new(Part)
	prev = head

	i = 0
	for i < len(bytes) && bytes[i] != 0 {
		part, j, err = parsePart(bytes[i:])
		if err != nil {
			return nil, i + j, err
		}

		i += j + 1

		prev.Next = part
		prev = part
	}
	
	return head.Next, i, nil
}

func ParseSecstore(bytes []byte) (*Secstore, error) {
	var err error
	var secstore *Secstore

	secstore = new(Secstore)

	secstore.Parts = new(Part)
	secstore.Parts.Name = "/"
	
	secstore.Parts.SubParts, _, err = parseParts(bytes)
	if err != nil {
		return nil, err
	}
	
	return secstore, nil
}

func (part *Part) ToBytes() []byte {
	var bytes []byte = []byte(nil)

	bytes = append(bytes, []byte(part.Name)...)
	bytes = append(bytes, 0)

	if part.Data == "" {
		for part := part.SubParts; part != nil; part = part.Next {
			bytes = append(bytes, part.ToBytes()...)
		}
	} else {
		bytes = append(bytes, []byte(part.Data)...)
	}
	
	bytes = append(bytes, 0)

	return bytes
}

func (store *Secstore) ToBytes() []byte {
	var bytes []byte = []byte(nil)

	for part := store.Parts.SubParts; part != nil; part = part.Next {
		bytes = append(bytes, part.ToBytes()...)
	}

	return bytes
}

func (part *Part) findSubPart(name string) *Part {
	for p := part.SubParts; p != nil; p = p.Next {
		if name == p.Name {
			return p
		}
	}

	return nil
}

func (store *Secstore) FindPart(name string) *Part {
	return store.Parts.findSubPart(name)
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

func (part *Part) Print() {
	if part.Data == "" {
		for p := part.SubParts; p != nil; p = p.Next {
			fmt.Printf("%s", p.Name)
			if p.Data == "" {
				fmt.Printf("/")
			}
			fmt.Printf("\n")
		}
	} else {
		fmt.Println(part.Data)
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

func (part *Part) addPart(p *Part) {
	p.Next = part.SubParts 
	part.SubParts = p
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
		store.Parts.addPart(part)
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

func (store *Secstore) MakeNewDirPart(name string) {
	var part *Part

	part = store.FindPart(name)
	if part != nil {
		fmt.Println(name, "already exists. Not adding.")
		return
	}

	part = new(Part)

	part.Data = ""
	part.SubParts = nil
	fmt.Println("Adding directory with name: ", name)
	store.Parts.addPart(part)
}