package main

import (
	"fmt"
	"os"
	"unicode"
	"strings"
	"github.com/peterh/liner"
)

func RunRepl(store *Secstore) {
	var sections []string
	var line string
	var err error

	liner := liner.NewLiner()
	defer liner.Close()

	for {
		line, err = liner.Prompt("> ")
		if err != nil {
			break
		}

		liner.AppendHistory(line)

		sections = splitSections(line)

		if strings.HasPrefix("quit", sections[0]) {
			break
		} else {
			err = EvalCommand(store, sections)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			} 
		}
	}
	
	fmt.Fprintln(os.Stderr, "Exiting...")
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

func EvalCommand(store *Secstore, line []string) error {
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
