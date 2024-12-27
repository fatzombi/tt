package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type taskItem struct {
	title string
}

func (i taskItem) Title() string       { return i.title }
func (i taskItem) Description() string { return "" }
func (i taskItem) FilterValue() string { return i.title }

const (
	modeNormal = iota
	modeInput
)

type taskSelector struct {
	list      list.Model
	textInput textinput.Model
	choice    string
	quitting  bool
	mode      int
}

func getRecentTasks() []list.Item {
	sessions, err := loadSessions()
	if err != nil {
		return []list.Item{taskItem{title: "New Task"}}
	}

	seen := make(map[string]bool)
	var items []list.Item

	// Add "New Task" as the first option
	items = append(items, taskItem{title: "New Task"})
	seen["New Task"] = true

	// Get unique tasks from recent sessions
	for _, session := range sessions {
		if !seen[session.Task] {
			items = append(items, taskItem{title: session.Task})
			seen[session.Task] = true
		}
	}

	return items
}

func initialTaskSelector() taskSelector {
	items := getRecentTasks()

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a task"
	l.SetShowHelp(false)

	ti := textinput.New()
	ti.Placeholder = "Enter task name"
	ti.Focus()

	return taskSelector{
		list:      l,
		textInput: ti,
		mode:      modeNormal,
	}
}

func (m taskSelector) Init() tea.Cmd {
	return nil
}

func (m taskSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 2)
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeNormal:
			switch msg.Type {
			case tea.KeyEnter:
				m.choice = m.list.SelectedItem().(taskItem).Title()
				if m.choice == "New Task" {
					m.mode = modeInput
					return m, nil
				}
				return m, tea.Quit

			case tea.KeyCtrlC:
				m.quitting = true
				return m, tea.Quit
			}

		case modeInput:
			switch msg.Type {
			case tea.KeyEnter:
				input := strings.TrimSpace(m.textInput.Value())
				if input != "" {
					m.choice = input
					return m, tea.Quit
				}
			case tea.KeyEsc:
				m.mode = modeNormal
				return m, nil
			}

			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m taskSelector) View() string {
	if m.mode == modeInput {
		return fmt.Sprintf("\nEnter task name:\n\n%s\n\n(press esc to cancel)", m.textInput.View())
	}
	return "\n" + m.list.View()
}

func promptForTask() (string, error) {
	p := tea.NewProgram(initialTaskSelector())
	m, err := p.Run()
	if err != nil {
		return "", err
	}

	selector := m.(taskSelector)
	if selector.quitting {
		return "", fmt.Errorf("selection cancelled")
	}

	return selector.choice, nil
}
