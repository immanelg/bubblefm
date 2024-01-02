package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	editor      string
	opener      string
	dirsfirst   bool
	dirsonly    bool
	cyclescroll bool
	preview     bool
	sort        SortType
	showhidden  bool
	// keys map[string]string // TODO: key parser & action methods methods
	// colors struct{} // TODO
	// icons map[string]string // TODO
}

func syntaxErr(lineNr int, line []string, expl string) error {
	return fmt.Errorf("line %d: syntax error: %v: %v", lineNr, expl, line)
}

func (self *config) source(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lineNr := 0
	for scanner.Scan() {
		lineNr++
		line := scanner.Text()
		tokens := strings.Fields(line)

		if len(tokens) == 0 || tokens[0] == "#" {
			continue
		}
		// The parser is extremely dumb i couldnt care less
		if len(tokens) != 2 {
			return syntaxErr(lineNr, tokens, "expected key value pair")
		}
		switch tokens[0] {
		case "editor":
			self.editor = tokens[1]

		case "opener":
			self.opener = tokens[1]

		case "cyclescroll":
			cyclescroll, err := strconv.ParseBool(tokens[1])

			if err != nil {
				return syntaxErr(lineNr, tokens, "invalid boolean")
			}
			self.cyclescroll = cyclescroll

		case "dirsfirst":
			dirsfirst, err := strconv.ParseBool(tokens[1])

			if err != nil {
				return syntaxErr(lineNr, tokens, "invalid boolean")
			}
			self.dirsfirst = dirsfirst

		case "preview":
			preview, err := strconv.ParseBool(tokens[1])

			if err != nil {
				return syntaxErr(lineNr, tokens, "invalid boolean")
			}
			self.preview = preview

		case "dirsonly":
			dirsonly, err := strconv.ParseBool(tokens[1])

			if err != nil {
				return syntaxErr(lineNr, tokens, "invalid boolean")
			}
			self.dirsonly = dirsonly

		case "showhidden":
			nohidden, err := strconv.ParseBool(tokens[1])

			if err != nil {
				return syntaxErr(lineNr, tokens, "invalid boolean")
			}
			self.showhidden = nohidden

		case "sort":
			switch tokens[1] {
			case "name":
				self.sort = SortName
			case "modified":
				self.sort = SortModified
			case "size":
				self.sort = SortSize
			default:
				return syntaxErr(lineNr, tokens, "invalid sorting method")
			}

		default:
			return syntaxErr(lineNr, tokens, "invalid key")
		}
	}
	return nil
}

func envOr(variable string, def string) string {
	if v := os.Getenv(variable); v == "" {
		return def
	} else {
		return v
	}
}

func initConfig() (c config) {
	c.editor = envOr("EDITOR", "nvim")
	c.opener = envOr("OPENER", "xdg-open")
	c.cyclescroll = true
	c.dirsfirst = true
	c.dirsonly = false
	c.preview = true
	c.showhidden = false
	c.sort = SortName

	return
}
