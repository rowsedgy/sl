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
	name       string
	ip         string
	webip      string
	pubauth    bool
	user       string
	password   string
	key        string
	tunnel     bool
	tunnelHost string
	legacy     bool
}

type listKeyMap struct {
	toggleDetails key.Binding
	openWeb       key.Binding
}

// additional key info on helpbar
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

func (i Item) Title() string       { return string(i.name) }
func (i Item) Description() string { return string("IP: " + i.ip + " WEB: " + i.webip) }
func (i Item) FilterValue() string { return string(i.name) }
func (i Item) User() string        { return string(i.user) }

type model struct {
	list       list.Model
	keys       *listKeyMap
	showDetail bool
	startSSH   bool
}

const fileName = "sl-connections.json"

type cfg struct {
	filepath    string
	connections connections
}

func (c *cfg) initialModel() model {
	// set up keys
	listKeys := newListKeyMap()
	newList, err := c.generateList()
	if err != nil {
		log.Fatal(err)
	}
	newList.Title = "SL"
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
		// ignore keybinds if filtering
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.startSSH = true
			return m, tea.Quit
		case "i":
			m.showDetail = !m.showDetail
		case "w":
			go func() {
				if err := handleWebLink(m); err != nil {
					log.Println(err)
				}
			}()
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
			Render(fmt.Sprintf("Name: %s\nIP: %s\nUser: %s\nWeb IP: %s\nPubKey: %s\nKeyPath: %s\nTunnel: %s\nTunnelHost: %s\nLegacy: %s\n\nPress \"i\" to go back to list.",
				selected.name,
				selected.ip,
				selected.user,
				selected.webip,
				fmt.Sprintf("%v", selected.pubauth),
				selected.key,
				fmt.Sprintf("%v", selected.tunnel),
				selected.tunnelHost,
				fmt.Sprintf("%v", selected.legacy),
			))
		return box
	}

	return m.list.View()
}

func main() {
	fullFilePath, err := checkConfigFile(fileName)
	if err != nil {
		fmt.Println("ERROR -", err)
	}
	globalCFG := cfg{
		filepath: fullFilePath,
	}

	args := os.Args[1:]

	if len(args) != 0 {
		globalCFG.handleArgs(args)
	}

	p := tea.NewProgram(
		globalCFG.initialModel(),
		tea.WithAltScreen(),
	)

	m, err := p.Run()
	if err != nil {
		fmt.Println("Error starting program:", err)
		os.Exit(1)
	}

	err = globalCFG.handleSSHSession(m)
	if err != nil {
		log.Fatal(err)
	}
}

func checkConfigFile(file string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	filePath := homeDir + "/.config/" + file

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return filePath, nil
}
