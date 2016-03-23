package main

import (
	"fmt"
	"strings"
)

type Part struct {
	Name string

	Data string
	SubParts *Part

	Next *Part
}

func ParsePart(bytes []byte) (*Part, int, error) {
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
		part.SubParts, j, err = ParseParts(bytes[k+1:])
		k = k + 1 + j
		if err != nil {
			return nil, k, err
		}
	}

	return part, k, nil
}

func ParseParts(bytes []byte) (*Part, int, error) {
	var i, j int
	var part, head, prev *Part
	var err error

	head = new(Part)
	prev = head
	
	i = 0
	for i < len(bytes) && bytes[i] != 0 {
		part, j, err = ParsePart(bytes[i:])
		if err != nil {
			return nil, i + j, err
		}

		i += j + 1
		
		prev.Next = part
		prev = part
	}

	return head.Next, i, nil
}

func (part *Part) ToBytes() []byte {
	var bytes []byte = []byte(nil)

	bytes = append(bytes, []byte(part.Name)...)
	bytes = append(bytes, 0)

	if part.Data == "" {
		bytes = append(bytes, 0)
		for part := part.SubParts; part != nil; part = part.Next {
			bytes = append(bytes, part.ToBytes()...)
		}
	} else {
		bytes = append(bytes, []byte(part.Data)...)
	}
	
	bytes = append(bytes, 0)

	return bytes
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

func (part *Part) FindSub(pathRaw string) *Part {
	path := strings.SplitN(pathRaw, "/", 2)

	for p := part.SubParts; p != nil; p = p.Next {
		if p.Name == path[0] {
			if len(path) > 1 && len(path[1]) > 0 {
				return p.FindSub(path[1])
			} else {
				return p
			}
		}
	}

	return nil
}
