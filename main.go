package main

import (
	"brain/form"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var style = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder(), true, false)

type screenState uint

const (
	listView screenState = iota
	formView
)

type Model struct {
	list     list.Model
	err      error
	selected *Note
	focused  bool
	mdView   viewport.Model
	renderer *glamour.TermRenderer
	form     form.NewNoteFormModel
	screenState
}

func NewModel() (*Model, error) {
	vp := viewport.New(100, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100))

	newNoteForm := form.InitialNewNoteFormModel()

	return &Model{
		renderer:    renderer,
		mdView:      vp,
		form:        newNoteForm,
		screenState: listView,
	}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}, width, height)
	m.list.Title = "Notes"
	m.list.SetItems([]list.Item{
		Note{title: "Day1", description: "holahola", content: `# New day same shit`, createdAt: time.Now()},
		Note{title: "Day1", description: "holahola", content: `# New day same shit`, createdAt: time.Now()},
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.screenState == formView {
		return m.form.Update(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.initList(msg.Width, msg.Height)
		return m, cmd
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", " ":
			if m.screenState == listView {
				if m.selected != nil {
					m.selected = nil
				} else {
					selectedNote, _ := m.list.SelectedItem().(Note)
					m.selected = &selectedNote
					out, _ := m.renderer.Render(m.selected.content)
					m.mdView.SetContent(out)
				}
				m.mdView, cmd = m.mdView.Update(msg)
				cmds = append(cmds, cmd)
			}
		case "n":
			m.screenState = formView
		}
	}

	if m.screenState == listView {
		m.list, cmd = m.list.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.screenState == formView {
		return m.form.View()
	}
	if m.selected != nil {
		return m.mdView.View()
	}
	return style.Render(m.list.View())
}

func main() {
	m, err := NewModel()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := tea.NewProgram(m)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
