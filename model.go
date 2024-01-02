package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type SortType uint

const (
	SortName SortType = iota
	SortModified
	SortSize
)

type model struct {
	files  []File
	cursor int
	empty  bool
	cwd    string

	width, height int
	topIndex      int

	preview []string

	selections set[string]
	bookmarks  map[string]string

	config config

	currentView ViewType
	status      status
}

func defaultBookmarks() map[string]string {
	// TODO: save bookmarks, selections to config file; also, dump all runtime settings to it.
	// TODO: mark for previous path
	return map[string]string{
		"1": expandHome("~"),
		"2": expandHome("~/Pictures"),
		"3": expandHome("~/Videos"),
		"4": expandHome("~/Downloads"),
		"5": "/",
	}
}

// height of normal view without statuslines/margins/paddings/borders
func (self *model) normalHeight() int {
	return self.height - 2
}

// height of normal view without margins/paddings/borders
func (self *model) normalWidth() int {
	return self.height - 2
}

func (self *model) bottomIndex() int {
	return self.topIndex + self.normalHeight() - 1
}

func (self *model) scroll(i int) {
	self.topIndex += i
}

func (self *model) moveCursor(i int) {
	self.cursor += i
	self.syncCursor()
}

func (self *model) toggleHidden() {
	self.config.showhidden = !self.config.showhidden
	self.syncCursor()
}

func (self *model) togglePreview() {
	self.config.preview = !self.config.preview
}

func (self *model) toggleDirsonly() {
	self.config.dirsonly = !self.config.dirsonly
	self.syncCursor()
}

// files that are visible in ui
func (self *model) visibleFiles() (files []File) {
	files = self.files
	if !self.config.showhidden {
		files = filter(&files, func(f File) bool {
			return !strings.HasPrefix(f.Name, ".")
		})
	}
	if self.config.dirsonly {
		files = filter(&files, func(f File) bool {
			return f.IsDir
		})
	}
	return files
}

// make cursor be in the bounds of model.files indices
func (self *model) syncCursor() {
	if self.empty {
		self.cursor = 0
		return
	}
	l := self.len()
	if self.cursor > l-1 {
		self.cursor = l - 1
	} else if self.cursor < 0 {
		self.cursor = 0
	}
}

// make ui bounds follow cursor
func (self *model) syncBounds() {
	if self.cursor < self.topIndex {
		self.topIndex = self.cursor
	} else if self.cursor > self.bottomIndex() {
		for self.bottomIndex() < self.cursor {
			self.topIndex++
		}
	}
}

// length of self.files that should be visible to ui based on config
func (self *model) len() int {
	return len(self.visibleFiles())
}

func (self *model) sortFiles() {
	sort.SliceStable(self.files, func(i, j int) bool {
		switch self.config.sort {
		case SortName:
			return self.files[i].Name < self.files[j].Name
		case SortModified:
			return self.files[i].Modified.Before(self.files[j].Modified)
		case SortSize:
			return self.files[i].Size < self.files[j].Size
		default:
			panic("unknown SortType")
		}
	})
	if self.config.dirsfirst {
		sort.SliceStable(self.files, func(i, j int) bool {
			if self.files[i].IsDir && !self.files[j].IsDir {
				return true
			}
			return false
		})
	}
}

func (self model) current() File {
	// TODO: return (ok bool, File)?
	if self.empty {
		panic("model.current is called but model.empty == true")
	}
	return self.visibleFiles()[self.cursor]
}

func (self model) open() (model, tea.Cmd) {
	if self.empty {
		return self, nil
	}
	current := self.current()
	if current.IsDir {
		return self, refreshFiles(current.Path)
	} else {
		return self, openCmd(self.config.editor, current.Name)
	}
}

func (self model) sortby(s SortType) {
	switch s {
	case SortName:
		self.status = newStatus("sort by time", false)
	case SortModified:
		self.status = newStatus("sort by time", false)
	case SortSize:
		self.status = newStatus("sort by time", false)
	default:
		panic("unknown SortType")
	}
	self.config.sort = s
	self.sortFiles()
}

func (self *model) refreshPreview() (model, tea.Cmd) {
	self.preview = []string{"..."}
	var cmd tea.Cmd
	if !self.empty {
		cmd = previewCmd(self.current().Path)
	}
	return *self, cmd
}

func (self *model) onKey(key string) (tea.Model, tea.Cmd) {
	if self.currentView == ViewHelp {
		self.currentView = ViewFiles
		return self, nil
	}

	if self.currentView == ViewSelections {
		self.currentView = ViewFiles
		return self, nil
	}

	switch key {

	case "j", "down":
		self.moveCursor(1)
		self.syncBounds()
		return self.refreshPreview()

	case "k", "up":
		self.moveCursor(-1)
		self.syncBounds()
		return self.refreshPreview()

	case "ctrl+d", "pgdown":
		self.moveCursor((self.bottomIndex() - self.topIndex) / 2)
		self.syncBounds()
		return self.refreshPreview()

	case "ctrl+u", "pgup":
		self.moveCursor(-(self.bottomIndex() - self.topIndex) / 2)
		self.syncBounds()
		return self.refreshPreview()

	case "g", "home":
		self.moveCursor(-self.len())
		self.syncBounds()
		return self.refreshPreview()

	case "G", "end":
		self.moveCursor(self.len())
		self.syncBounds()
		return self.refreshPreview()

	case "n":
		self.sortby(SortName)
		return self.refreshPreview()

	case "t":
		self.sortby(SortModified)
		return self.refreshPreview()

	case "s":
		self.sortby(SortSize)
		return self.refreshPreview()

	case "ctrl+l":
		return self, refreshFiles(self.cwd)

	case "v":
		// TODO: save selections to registers?
		if self.empty {
			break
		}
		self.selections.Toggle(self.current().Path)
		self.moveCursor(1)
		return self.refreshPreview()

	case "V":
		if self.empty {
			break
		}
		self.selections.Toggle(self.current().Path)
		self.moveCursor(-1)
		return self.refreshPreview()

	case "esc":
		self.selections.Clear()

		// TODO: yank paths to clipboard
		// TODO: marks like in Vim (m + letter, ' + letter, pasting to mark, etc)
	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		self.bookmarks[key] = self.cwd

	case "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9":
		key := string(key[1]) // second byte is the second grapheme in this case
		path, exists := self.bookmarks[key]
		if exists {
			return self, refreshFiles(path)
		}

		// TODO: stuff with symlinks
		// TODO: paste as symbolic links

	case "p":
		return self, tea.Sequence(copyFilesCmd(self.selections, self.cwd), refreshFiles(self.cwd))

	case "P":
		return self, tea.Sequence(moveFilesCmd(self.selections, self.cwd), refreshFiles(self.cwd))

	case "D":
		// TODO: prompt before deletion
		return self, tea.Sequence(deleteFilesCmd(self.selections), refreshFiles(self.cwd))

	case "h", "left":
		return self, refreshFiles(filepath.Dir(self.cwd))

	case "l", "right":
		return self.open()

	case "o":
		if self.empty {
			break
		}
		current := self.current()
		return self, openExternalCmd(self.config.opener, current.Name)

	case ".":
		self.toggleHidden()
		return self.refreshPreview()

	case "/":
		self.toggleDirsonly()
		return self.refreshPreview()

	case "f":
		self.togglePreview()

	case "q", "ctrl+c":
		return self, tea.Quit

	case "?":
		self.currentView = ViewHelp
		return self, nil

	case " ":
		self.currentView = ViewSelections
		return self, nil

	default:
		self.status = newStatus(fmt.Sprintf("unmapped key: %s", key), true)
		return self, clearStatusCmd(self.status.id)
	}
	return self, nil
}

func (self model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		return self.onKey(key)

	case tea.WindowSizeMsg:
		self.width, self.height = msg.Width, msg.Height

	case filesRefreshMsg:
		if msg.err != nil {
			self.status = newStatus(msg.err.Error(), true)
			return self, clearStatusCmd(self.status.id)
		}

		self.cwd = msg.cwd
		os.Chdir(self.cwd)
		self.files = msg.files
		self.sortFiles()
		self.empty = self.len() == 0
		self.cursor = 0
		self.topIndex = 0
		return self.refreshPreview()

	case copyFilesMsg:
		if msg.err != nil {
			self.status = newStatus(msg.err.Error(), true)
			return self, clearStatusCmd(self.status.id)
		}
		self.selections.Clear()

	case moveFilesMsg:
		if msg.err != nil {
			self.status = newStatus(msg.err.Error(), true)
			return self, clearStatusCmd(self.status.id)
		}
		self.selections.Clear()

	case deleteFilesMsg:
		if msg.err != nil {
			self.status = newStatus(msg.err.Error(), true)
			return self, clearStatusCmd(self.status.id)
		}
		self.selections.Clear()

	case clearStatusMsg:
		if self.status.id == msg.id {
			self.status = newStatus("", false)
		}
		return self, nil

	case processFininishedMsg:
		if msg.err != nil {
			self.status = newStatus(fmt.Sprintf("process finished: %v", msg.err.Error()), true)
			return self, clearStatusCmd(self.status.id)
		}

	case previewMsg:
		if msg.err != nil {
			self.status = newStatus(fmt.Sprintf("previewing %v: %v", msg.path, msg.err.Error()), true)
			return self, clearStatusCmd(self.status.id)
		}
		if !self.empty && msg.path == self.current().Path {
			self.preview = msg.lines
		}
	}
	return self, nil
}

func (self model) Init() tea.Cmd {
	return refreshFiles(self.cwd)
}

func newModel(cwd string, config config) model {
	return model{
		files:       []File{},
		cwd:         cwd,
		cursor:      0,
		empty:       true,
		width:       0,
		height:      0,
		topIndex:    0,
		selections:  make(set[string]),
		bookmarks:   defaultBookmarks(),
		currentView: ViewFiles,
		config:      config,
		status:      status{},
	}
}
