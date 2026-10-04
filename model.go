package main

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textInput textinput.Model // items on the to-do list
	help      help.Model
	keys      keyMap
	quitting  bool
	repos     []repo
}

func initialModel() tea.Model {
	m := &model{
		textInput: textinput.New(),
		help:      help.New(),
		keys:      keys,
		quitting:  false,
	}

	m.textInput.Prompt = "Charm/"
	m.textInput.Placeholder = "Enter repo name"
	m.textInput.Focus()
	m.textInput.CharLimit = 100
	m.textInput.ShowSuggestions = true
	m.textInput.SetVirtualCursor(false)

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getRepos, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.textInput.SetWidth(msg.Width / 3)

	case repoGetSuccess:
		var suggestions []string
		for _, r := range msg {
			suggestions = append(suggestions, r.Name)
		}
		m.textInput.SetSuggestions(suggestions)
		return m, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() tea.View {

	if m.quitting {
		return tea.NewView("Terminating program, Bye !!\n")
	}

	helpView := m.help.View(m.keys)

	header := headerStyle.Render("Type th repo name, we will complet it ^^!")
	m.textInput.Prompt = promptStyle.Render(m.textInput.Prompt)
	t := textInputStyle.Render(m.textInput.View())

	v := lipgloss.JoinVertical(lipgloss.Left, header, t, helpView)

	view := tea.NewView(v)
	view.AltScreen = true
	return view
}

var (
	textInputStyle = lipgloss.NewStyle().MarginLeft(10)
	promptStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	headerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true).Underline(true)
)
