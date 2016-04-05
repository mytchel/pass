package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/peterh/liner"
)

func RunRepl(secstore *Secstore) {
	var sections []string
	var line string
	var err error
	var quit bool = false

	liner := liner.NewLiner()
	defer liner.Close()

	for {
		line, err = liner.Prompt("> ")
		if err != nil {
			break
		}

		liner.AppendHistory(line)

		sections = splitSections(line)

		quit, err = evalLine(secstore, sections)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if quit {
			break
		}
	}
}

func splitSections(s string) (sections []string) {
	var i, j int
	var quote bool = false
	var section string
	
	i = 0
	for i < len(s) {
		section = ""
		for j = i; j < len(s); j++ {
			if s[j] == '\'' {
				quote = !quote
			} else if unicode.IsSpace(rune(s[j])) && !quote {
				break		
			} else {
				section = section + string(s[j])
			}
		}

		sections = append(sections, section)

		for i = j; i < len(s); i++ {
			if !unicode.IsSpace(rune(s[i])) {
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
	case "c":
		if len(line) < 2 {
			return false, fmt.Errorf("Usage: c 'dir'\n\tUsing .. as dir will go to the parent directory.")
		} else {
			secstore.ChangeDir(line[1])
		}
	default:
		return false, fmt.Errorf("%s: not a command", line[0])
	}

	return false, nil
}
