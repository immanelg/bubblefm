package main

import (
	"fmt"
	"strings"

	lipgloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type ViewType byte

const (
	ViewFiles ViewType = iota
	ViewHelp
	ViewSelections
)

func (self *model) helpView() string {
	tbl := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#18ff20"))).
		Headers("Key", "Action").
		Width(self.width)

	tbl.Row("q", "Quit application")
	tbl.Row("?", "Open this help")
	tbl.Row("j/k/g/G", "Down/Up")
	tbl.Row("s/t/n", "Sort by size/time/name")
	tbl.Row("h/l", "Updir/Downdir")
	tbl.Row("v/V", "Select file")
	tbl.Row("p", "Copy selections")
	tbl.Row("P", "Move selections")
	tbl.Row("D", "Remove selections")
	tbl.Row("esc", "Clear selections")
	tbl.Row("o", "Open in app")
	tbl.Row(".", "Toggle hidden")
	tbl.Row("f", "Toggle preview")
	tbl.Row("C-l", "Reload files")
	tbl.Row("{1..9}", "Bookmark this dir")
	tbl.Row("f{1..9}", "Go to bookmark")

	return lipgloss.JoinVertical(lipgloss.Left, tbl.Render(), "Press any key to close help")
}

func (self *model) toplineView() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#006cb6")).Bold(true)
	var d string
	if self.cwd == "/" {
		d = "/"
	} else {
		d = withTilde(self.cwd) + "/"
	}
	dirname := style.Render(d)

	var basename string
	if !self.empty {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#daf52e"))
		f := self.current().Name
		basename = style.Render(f)
	}

	view := dirname + basename
	return view
}

func (self *model) statusView() (view string) {
	if self.status.isErr {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
		view = style.Render(fmt.Sprintf("error: %v", self.status.text))
	} else {
		style := lipgloss.NewStyle().Italic(true)
		view = style.Render(self.status.text)
	}
	return
}

func (self *model) fileListView() string {
	items := []string{}
	files := self.visibleFiles()

	begin := self.topIndex
	end := min(self.bottomIndex(), len(files)-1)
	if begin > end {
		panic(fmt.Sprintf("top > bottomOfFiles: %v > %v", begin, end))
	}

	for i := begin; i <= end; i++ {
		file := files[i]

		styleName := lipgloss.NewStyle()

		var fileIcon string
		if file.IsDir {
			styleName = styleName.Foreground(lipgloss.Color("#3071ff"))
			fileIcon = ""
		} else {
			styleName = styleName.Foreground(lipgloss.Color("#ffffff"))
			fileIcon = ""
		}

		var selectedIcon string
		if self.selections.Contains(file.Path) {
			selectedIcon = "*"
			styleName = styleName.Bold(true)
		} else {
			selectedIcon = " "
		}

		mode := file.Mode.String() // TODO: prettify
		styleMode := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffd35e"))
		viewMode := styleMode.Render(mode)

		styleModified := lipgloss.NewStyle().Foreground(lipgloss.Color("#bbbbbb"))
		viewModified := styleModified.Render(file.Modified.Format("2006-01-02 15:04"))

		if i == self.cursor {
			styleName = styleName.Background(lipgloss.Color("#616161"))
		}

		viewFilename := styleName.Render(fmt.Sprintf("%s %s  %s", selectedIcon, fileIcon, file.Name))

		itemView := viewFilename
		width := 45

		if self.config.preview && self.width/2 >= width || !self.config.preview && self.width >= width {
			itemMetadataView := fmt.Sprintf("%s %s", viewMode, viewModified)
			itemView = itemMetadataView + itemView
		}

		if remain := self.width/2 - lipgloss.Width(itemView); remain > 0 {
			itemView += strings.Repeat(" ", remain)
		}
		items = append(items, itemView)
	}
	remain := self.normalHeight() - len(items)
	for i := 0; i < remain; i++ {
		items = append(items, " ")
	}

	style := lipgloss.NewStyle().MaxHeight(self.normalHeight())
	view := style.Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	return view
}

func (self *model) previewView() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#afafaf")).MaxHeight(self.normalHeight())
	view := style.Render(lipgloss.JoinVertical(lipgloss.Left, self.preview...))
	return view
}

func (self *model) selectionsView() string {
	var lines []string
	for line := range self.selections {
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		lines = []string{"no selections"}
	}
	view := lipgloss.JoinVertical(lipgloss.Left,
		lines...,
	)
	return view
}

func (self model) View() string {
	if self.height < 3 || self.width < 10 {
		return "..."
	}

	var mainView string

	if self.currentView == ViewHelp {
		mainView = self.helpView()
	} else if self.currentView == ViewSelections {
		mainView = self.selectionsView()
	} else if self.empty {
		mainView = "very empty here, innit?"
	} else if self.config.preview {
		mainView = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().MaxWidth(self.width/2).Render(self.fileListView()),
			lipgloss.NewStyle().MaxWidth(self.width/2).Render(self.previewView()),
		)
	} else {
		mainView = self.fileListView()
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		self.toplineView(),
		mainView,
		self.statusView(),
	)
}
