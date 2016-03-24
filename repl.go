package main

import (
	"fmt"
	"os"
	"unicode"
)

const (
	LineLen = 512
)

func RunRepl(secstore *Secstore) {
	var bytes []byte = make([]byte, LineLen)
	var sections []string
	var line []byte
	var n int
	var err error
	var quit bool = false

	for {
		fmt.Printf("> ")
		n, err = os.Stdin.Read(bytes)
		if err != nil {
			break
		}
	
		for i := 0; i < n; i++ {
			if bytes[i] == '\n' {
				line = bytes[:i]
				break
			}
		}

		sections = splitSections(line)

		quit, err = evalLine(secstore, sections)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if quit {
			break
		}
	}
}

func splitSections(bytes []byte) []string {
	var s, i, j int
	var quote bool = false
	var sections []string = []string(nil)
	var section []byte = make([]byte, LineLen)
	
	i = 0
	for i < len(bytes) {
		s = 0
		for j = i; j < len(bytes); j++ {
			if bytes[j] == '\'' {
				quote = !quote
			} else if unicode.IsSpace(rune(bytes[j])) && !quote {
				break		
			} else {
				section[s] = bytes[j]
				s++
			}
		}

		sections = append(sections, string(section[:s]))

		for i = j; i < len(bytes); i++ {
			if !unicode.IsSpace(rune(bytes[i])) {
				break		
			}
		}

}
	
	return sections
}

func evalLine(secstore *Secstore, line []string) (bool, error) {
	if len(line) < 1 {
		return false, nil
	}

	switch (line[0]) {
	case "q":
		return true, nil
	case "a":
		if len(line) != 2 {
			return false, fmt.Errorf("Usage: a 'new pass name'")
		} else {
			secstore.MakeNewPart(line[1])
		}
	case "m":
		if len(line) != 2 {
			return false, fmt.Errorf("Usage: m 'new dir name'")
		} else {
			secstore.MakeNewDirPart(line[1])
		}
	case "s":
		if len(line) == 1 {
			secstore.List()
		} else {
			for i := 1; i < len(line); i++ {
				secstore.ShowPart(line[i])
			}
		}
	case "d": 
		if len(line) < 2 {
			return false, fmt.Errorf("Usage: d 'pass or dir name'")
		} else {
			secstore.RemovePart(line[1])
		}
	case "e":
		if len(line) < 2 {
			return false, fmt.Errorf("Usage: e 'pass name'")
		} else {
			secstore.EditPart(line[1])
		}
	default:
		return false, fmt.Errorf("%s: not a command", line[0])
	}

	return false, nil
}
