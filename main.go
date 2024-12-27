package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
)

const (
	defaultDuration = 25 * time.Minute
)

type model struct {
	task       string
	duration   time.Duration
	startTime  time.Time
	elapsed    time.Duration
	progress   progress.Model
	width      int
	shouldSave bool
}

func initialModel(task string, duration time.Duration) model {
	p := progress.New(progress.WithDefaultGradient())
	return model{
		task:       task,
		duration:   duration,
		startTime:  time.Now(),
		progress:   p,
		shouldSave: true,
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = msg.Width - 4
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			if m.shouldSave {
				actualDuration := time.Since(m.startTime)
				saveSession(m.task, actualDuration, m.startTime)
				fmt.Printf("\nSaved partial session: %s (%.2f seconds)\n",
					m.task,
					actualDuration.Seconds(),
				)
			}
			return m, tea.Quit
		}

	case tickMsg:
		m.elapsed = time.Since(m.startTime)
		if m.elapsed >= m.duration {
			if m.shouldSave {
				saveSession(m.task, m.duration, m.startTime)
				m.shouldSave = false
				notify(m.task)
			}
			return m, tea.Quit
		}
		return m, tick()
	}

	return m, nil
}

func (m model) View() string {
	percent := float64(m.elapsed) / float64(m.duration)
	remainingTime := m.duration - m.elapsed

	return fmt.Sprintf(
		"\n  %s: %s remaining\n\n%s\n",
		m.task,
		remainingTime.Round(time.Second),
		m.progress.ViewAs(percent),
	)
}

func notify(task string) {
	// Play the Crystal sound
	go exec.Command("afplay", "/System/Library/Sounds/Crystal.aiff").Run()

	// Show notification
	beeep.Alert("Work Timer", fmt.Sprintf("Finished working on: %s", task), "")
}

func main() {
	var task string
	var duration time.Duration = defaultDuration

	// Handle summary command first
	if len(os.Args) > 1 && os.Args[1] == "summary" {
		printDailySummary()
		return
	}

	// Parse duration if it's the first argument
	if len(os.Args) > 1 {
		if minutes, err := strconv.Atoi(os.Args[1]); err == nil {
			duration = time.Duration(minutes) * time.Minute
			// If duration was first arg, and no task specified, prompt for task
			if len(os.Args) == 2 {
				selectedTask, err := promptForTask()
				if err != nil {
					fmt.Println("Error:", err)
					os.Exit(1)
				}
				if selectedTask == "" {
					fmt.Println("No task selected")
					os.Exit(1)
				}
				task = selectedTask
			} else {
				task = os.Args[2] // Task was provided after duration
			}
		} else {
			// First arg wasn't a number, so it must be the task name
			task = os.Args[1]
			// Check if second arg is duration
			if len(os.Args) > 2 {
				if minutes, err := strconv.Atoi(os.Args[2]); err == nil {
					duration = time.Duration(minutes) * time.Minute
				}
			}
		}
	} else {
		// No args provided, prompt for task
		selectedTask, err := promptForTask()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		if selectedTask == "" {
			fmt.Println("No task selected")
			os.Exit(1)
		}
		task = selectedTask
	}

	p := tea.NewProgram(initialModel(task, duration))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
