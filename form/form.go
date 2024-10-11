package form

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

const (
	title = iota
	content
)

type errMsg error

type NewNoteFormModel struct {
	err     error
	inputs  []textinput.Model
	focused int
}

func InitialNewNoteFormModel() NewNoteFormModel {
	inputs := make([]textinput.Model, 2)
	inputs[title] = textinput.New()
	inputs[title].Placeholder = "My new fantastic note"
	inputs[title].Focus()
	inputs[title].CharLimit = 120
	inputs[title].Width = 50
	inputs[title].Prompt = ""

	inputs[content] = textinput.New()
	inputs[content].Placeholder = `# The title of your fantastic note`
	inputs[content].Focus()
	inputs[content].Width = 50
	inputs[content].Prompt = ""

	return NewNoteFormModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (nf NewNoteFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (nf NewNoteFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(nf.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if nf.focused == len(nf.inputs)-1 {
				return nf, tea.Quit
			}
			nf.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return nf, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			nf.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			nf.nextInput()
		}
		for i := range nf.inputs {
			nf.inputs[i].Blur()
		}
		nf.inputs[nf.focused].Focus()
	case errMsg:
		nf.err = msg
		return nf, nil
	}
	for i := range nf.inputs {
		nf.inputs[i], cmds[i] = nf.inputs[i].Update(msg)
	}

	return nf, tea.Batch(cmds...)
}

func (nf NewNoteFormModel) View() string {
	return fmt.Sprintf(
		` 
 %s
 %s

 %s 
 %s

 %s
`,
		inputStyle.Width(30).Render("Title"),
		nf.inputs[title].View(),
		inputStyle.Width(30).Render("Content"),
		nf.inputs[content].View(),
		continueStyle.Render("Create"),
	) + "\n"
}

func (nf *NewNoteFormModel) nextInput() {
	nf.focused = (nf.focused + 1) % len(nf.inputs)
}

func (nf *NewNoteFormModel) prevInput() {
	nf.focused--
	if nf.focused < 0 {
		nf.focused = len(nf.inputs) - 1
	}
}
