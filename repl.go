package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/peterh/liner"
)

func RunRepl(store *Secstore, line *liner.State) {
	var sections []string
	var l string
	var err error

	line.SetCompleter(completer)

	for {
		l, err = line.Prompt("> ")
		if err != nil {
			break
		}

		line.AppendHistory(l)

		sections = splitSections(l)

		if len(sections) == 0 {

		} else if strings.HasPrefix("quit", sections[0]) {
			break
		} else {
			err = EvalCommand(store, line, sections)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
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

func EvalCommand(store *Secstore, line *liner.State, sections []string) error {
	args := []string(nil)
	if len(sections) > 1 {
		args = sections[1:]
	}

	if f, err := MatchCommand(sections[0]); err == nil {
		return f(store, line, args)
	} else {
		return err
	}
}

func completer(line string) []string {
	sections := splitSections(line)

	if len(sections) > 1 {
		return []string(nil)
	} else if len(sections) == 1 {
		var matches []string = []string(nil)

		for c, _ := range Commands {
			if strings.HasPrefix(c, sections[0]) {
				matches = append(matches, c)
			}
		}
		return matches
	} else {
		return []string(nil)
	}
}
