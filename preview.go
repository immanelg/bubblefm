package main

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type previewMsg struct {
	path  string
	lines []string
	err   error
}

// TODO: binary files, syntax highlighting, images, ...
func previewCmd(path string) tea.Cmd {
	return func() tea.Msg {
		height := 100

		var lines []string

		info, err := os.Stat(path)
		if err != nil {
			return previewMsg{err: err}
		}
		if info.IsDir() {
			lines = append(lines, fmt.Sprintf("directory %v:", path))

			files, err := os.ReadDir(path)
			if err != nil {
				return previewMsg{err: err}
			}
			for i, file := range files {
				if i > height {
					break
				}
				lines = append(lines, file.Name())
			}
		}

		file, err := os.Open(path)
		if err != nil {
			return previewMsg{err: err}
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for i := 0; scanner.Scan() && i < height; i++ {
			l := scanner.Text()
			lines = append(lines, l)
		}

		return previewMsg{path, lines, err}
	}
}
