package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/peterh/liner"
)

var lineReader *liner.State

var commands = map[string](func(*Secstore, []string) error) {
	"chpass": 	ChangePass,
	"add":		AddDataPart,
	"mkdir":	AddDirPart,
	"show":		ShowPart,
	"ls":		ShowPart,
	"rm":		RemovePart,
	"edit":		EditPart,
	"cd":		ChangeDir,
	"mv":		MovePart,
}

func RunRepl(store *Secstore) {
	var sections []string
	var line string
	var err error

	lineReader = liner.NewLiner()

	for {
		line, err = lineReader.Prompt("> ")
		if err != nil {
			break
		}

		lineReader.AppendHistory(line)

		sections = splitSections(line)

		err = evalLine(store, sections)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
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

func quit(store *Secstore) error {
	if err := SaveSecstore(store); err != nil {
		return err
	} else {
		lineReader.Close()
		os.Exit(0)
		return nil
	}
}

func evalLine(store *Secstore, line []string) error {
	if len(line) < 1 {
		return nil
	}

	if line[0] == "quit" {
		quit(store)
		return nil
	} else {
		if f := commands[line[0]]; f != nil {
			if len(line) > 1 {
				return f(store, line[1:])
			} else {
				return f(store, []string(nil))
			}
		} else {
			return fmt.Errorf("No command matching '%s' found.", line[0])
		}
	}
}

