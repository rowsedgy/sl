package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	name     string
	ip       string
	user     string
	password string
}

// items

func (i Item) Title() string       { return string(i.name) }
func (i Item) Description() string { return string(i.ip) }
func (i Item) FilterValue() string { return string(i.name) }
func (i Item) User() string        { return string(i.user) }

// func (i Item) Key() string         { return string(i.key) }

type model struct {
	list       list.Model
	showDetail bool
}

const filepath = "connections.json"

func initialModel() model {
	newList, err := generateList(filepath)
	if err != nil {
		log.Fatal(err)
	}
	newList.Title = "Lista de Nagios"
	newList.SetShowStatusBar(true)
	newList.SetFilteringEnabled(true)
	newList.SetShowPagination(true)
	newList.SetShowHelp(true)

	newList.KeyMap.Filter.SetKeys("/")
	newList.KeyMap.ClearFilter.SetKeys("esc")

	return model{list: newList}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			selected := m.list.SelectedItem().(Item)
			err := spawnSSHSession(
				selected.user,
				selected.password,
				selected.ip,
			)
			if err != nil {
				log.Fatal(err)
			}
		case "i":
			m.showDetail = !m.showDetail
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.showDetail {
		selected, ok := m.list.SelectedItem().(Item)
		if !ok {
			return "No item selected"
		}

		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(1, 2).
			Width(50).
			Render(fmt.Sprintf("Name: %s\nIP: %s\nUser: %s\nKey path: %s\n\nPress i to go back to list.",
				selected.name,
				selected.ip,
				selected.user,
				selected.password,
			))
		return box
	}

	return m.list.View()
}

func main() {
	generateList("connections.json")
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if err := p.Start(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
