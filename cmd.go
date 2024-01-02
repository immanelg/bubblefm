package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type filesRefreshMsg struct {
	files []File
	cwd   string
	err   error
}

type copyFilesMsg struct {
	err error
}

type deleteFilesMsg struct {
	err error
}

type moveFilesMsg struct {
	err error
}

type createFileMsg struct {
	err error
}

type clearStatusMsg struct {
	id int
}

type processFininishedMsg struct {
	err error
}

func clearStatusCmd(id int) tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return clearStatusMsg{id}
	})
}

// TODO: 1) messages for file operations progress
// 2) quit only after all file operation tasks are done & force quit

// TODO: cp files
func copyFilesCmd(selections set[string], toPath string) tea.Cmd {
	return func() tea.Msg {
		var msg copyFilesMsg
		// for selection := range selections {
		// 	src, err := os.Open(selection)
		// 	if err != nil {
		// 		msg.Err = err
		// 		return msg
		// 	}
		// 	defer src.Close()
		//
		// dest, err := os.Create(filepath.Join(path, selection))
		// 	if err != nil {
		// 		msg.Err = err
		// 		return msg
		// 	}
		// 	defer dest.Close()
		//
		// 	_, err = io.Copy(src, dest)
		// 	if err != nil {
		// 		msg.Err = err
		// 		return msg
		// 	}
		// }
		msg.err = errors.New("I don't feel like implementing it")
		return msg
	}
}

// TODO: mv files
func moveFilesCmd(selections set[string], toPath string) tea.Cmd {
	return func() tea.Msg {
		var msg moveFilesMsg
		msg.err = errors.New("I don't feel like implementing it")
		return msg
	}
}

// TODO: rm files
// trash?
func deleteFilesCmd(selections set[string]) tea.Cmd {
	return func() tea.Msg {
		var msg deleteFilesMsg
		msg.err = errors.New("I don't feel like implementing it")
		return msg
	}
}

// TODO: touch/mkdir
func createFileCmd(name string, dir string) tea.Cmd {
	return func() tea.Msg {
		var msg createFileMsg
		msg.err = errors.New("I don't feel like implementing it")
		return msg
	}
}

func refreshFiles(path string) tea.Cmd {
	return func() tea.Msg {
		var msg filesRefreshMsg

		var err error

		if err != nil {
			msg.err = err
			return msg
		}

		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			msg.err = fmt.Errorf("directory %v does not exists", path)
			return err
		}

		msg.cwd = path

		entries, err := os.ReadDir(path)
		if err != nil {
			msg.err = err
			return msg
		}

		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				msg.err = err
				return msg
			}
			file := File{
				Name:     info.Name(),
				Modified: info.ModTime(),
				Mode:     info.Mode(),
				Size:     info.Size(),
				IsDir:    info.IsDir(),
			}
			file.Path = filepath.Join(path, file.Name)
			msg.files = append(msg.files, file)
		}
		return msg
	}
}

func openCmd(program string, args ...string) tea.Cmd {
	cmd := tea.ExecProcess(exec.Command(program, args...), func(err error) tea.Msg {
		return processFininishedMsg{err}
	})
	return cmd
}

func openExternalCmd(program string, args ...string) tea.Cmd {
	return func() tea.Msg {
		err := exec.Command(program, args...).Start()
		return processFininishedMsg{err}
	}
}
