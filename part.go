package main

import (
	"fmt"
)

const (
	TypeDir  = 1
	TypeData = 2
)

type Part struct {
	Name string

	Type     int
	Data     string
	SubParts *Part

	Parent, Next *Part
}

func ParsePart(bytes []byte, parent *Part) (*Part, int, error) {
	var j, k int
	var part *Part
	var err error

	part = new(Part)
	part.Parent = parent

	for j = 0; j < len(bytes) && bytes[j] != 0; j++ {
	}
	if j == len(bytes) {
		return nil, j, fmt.Errorf("Error parsing part: Reached end.")
	}

	part.Name = string(bytes[:j])

	for k = j + 1; k < len(bytes) && bytes[k] != 0; k++ {
	}
	if k == len(bytes) {
		return nil, k, fmt.Errorf("Error parsing part: Reached end.")
	}

	/* Data part */
	if k > j+1 {
		part.Type = TypeData
		part.Data = string(bytes[j+1 : k])

		/* Sub tree */
	} else {
		part.Type = TypeDir
		part.SubParts, j, err = ParseParts(bytes[k+1:], part)
		k = k + 1 + j
		if err != nil {
			return nil, k, err
		}
	}

	return part, k, nil
}

func ParseParts(bytes []byte, parent *Part) (*Part, int, error) {
	var i, j int
	var part, head, prev *Part
	var err error

	head = new(Part)
	prev = head

	i = 0
	for i < len(bytes) && bytes[i] != 0 {
		part, j, err = ParsePart(bytes[i:], parent)
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

	if part.Type == TypeDir {
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
	if part.Type == TypeDir {
		if part.Parent != nil {
			fmt.Println("..")
		}

		for p := part.SubParts; p != nil; p = p.Next {
			fmt.Printf("%s", p.Name)
			if p.Type == TypeDir {
				fmt.Printf("/")
			}
			fmt.Printf("\n")
		}
	} else {
		fmt.Printf("%s", part.Data)
	}
}

func (part *Part) FindSub(path []string) *Part {
	if len(path) == 0 {
		return part
	} else if path[0] == "." || len(path[0]) == 0 {
		return part.FindSub(path[1:])
	} else if path[0] == ".." {
		if len(path) > 1 {
			return part.Parent.FindSub(path[1:])
		} else {
			return part.Parent
		}
	} else {
		for p := part.SubParts; p != nil; p = p.Next {
			if path[0] == p.Name {
				if len(path) > 1 {
					return p.FindSub(path[1:])
				} else {
					return p
				}
			}
		}
	}

	return nil
}

func (part *Part) AddPart(o *Part) error {
	var p *Part

	if part.Type != TypeDir {
		return fmt.Errorf("'%s' is not a directory.", part.Name)
	}

	if part.SubParts == nil {
		part.SubParts = o
	} else {
		for p = part.SubParts; p.Next != nil; p = p.Next {
		}
		p.Next = o
	}

	return nil
}

func (part *Part) RemovePart(o *Part) error {
	var p *Part

	if part.Type != TypeDir {
		return fmt.Errorf("'%s' is not a directory.", part.Name)
	}

	if part.SubParts == o {
		part.SubParts = part.SubParts.Next
	} else {
		for p = part.SubParts; p.Next != o; p = p.Next {
		}
		p.Next = o.Next
	}

	return nil
}
