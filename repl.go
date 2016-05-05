package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/peterh/liner"
)

func RunRepl(store *Secstore) {
	var sections []string
	var line string
	var err error

	LineReader = liner.NewLiner()

	for {
		line, err = LineReader.Prompt("> ")
		if err != nil {
			break
		}

		LineReader.AppendHistory(line)

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

func evalLine(store *Secstore, line []string) error {
	if len(line) < 1 {
		return nil
	}

	if f, err := MatchCommand(line[0]); err == nil {
		if len(line) > 1 {
			return f(store, line[1:])
		} else {
			return f(store, []string(nil))
		}
	} else {
		return err
	}
}
