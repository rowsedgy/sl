package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Item struct {
	name     string
	ip       string
	webip    string
	user     string
	password string
}

type listKeyMap struct {
	toggleDetails key.Binding
	openWeb       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		toggleDetails: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "toggle details"),
		),
		openWeb: key.NewBinding(
			key.WithKeys("w"),
			key.WithHelp("w", "open Web link"),
		),
	}
}

// items

func (i Item) Title() string       { return string(i.name) }
func (i Item) Description() string { return string(i.ip) }
func (i Item) FilterValue() string { return string(i.name) }
func (i Item) User() string        { return string(i.user) }

// func (i Item) Key() string         { return string(i.key) }

type model struct {
	list       list.Model
	keys       *listKeyMap
	showDetail bool
	startSSH   bool
}

const filepath = "connections.json"

func initialModel() model {
	// set up keys
	var listKeys = newListKeyMap()
	newList, err := generateList(filepath)
	if err != nil {
		log.Fatal(err)
	}
	newList.Title = "Lista de Nagios"
	newList.SetShowStatusBar(true)
	newList.SetFilteringEnabled(true)
	newList.SetShowPagination(true)
	newList.SetShowHelp(true)
	newList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleDetails,
			listKeys.openWeb,
		}
	}

	newList.KeyMap.Filter.SetKeys("/")
	newList.KeyMap.ClearFilter.SetKeys("esc")

	return model{list: newList, startSSH: false}
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
			m.startSSH = true
			return m, tea.Quit
		case "i":
			m.showDetail = !m.showDetail
		case "w":
			err := handleWebLink(m)
			if err != nil {
				log.Fatal(err)
			}
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

	m, err := p.Run()
	if err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}

	err = handleSSHSession(m)
	if err != nil {
		log.Fatal(err)
	}

}
