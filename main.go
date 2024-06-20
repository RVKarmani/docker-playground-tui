package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func initialModel() model {
	m := model{err: nil, dockerhubresponse: DockerHubResponse{}}

	m.dockeruserinput = textinput.New()
	m.dockeruserinput.Placeholder = "Enter docker username...."
	m.dockeruserinput.Focus()
	m.dockeruserinput.CharLimit = 156
	m.dockeruserinput.Width = 20

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case errMsg:
		m.dockerhubresponse = DockerHubResponse{}
		m.err = msg
		return m, nil

	case DockerHubResponse:
		m.err = nil
		m.dockerhubresponse = msg
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			return m, getDockerHubDetails(m.dockeruserinput.Value())
		}
	}

	m.dockeruserinput, cmd = m.dockeruserinput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	rightPaneContent := ""

	if m.err != nil {
		rightPaneContent = rightPaneContent + lipgloss.NewStyle().Foreground(hotPink).Render(m.err.Error())
	} else if (m.dockerhubresponse != DockerHubResponse{}) {
		pretty, _ := json.MarshalIndent(m.dockerhubresponse, "", "    ")
		rightPaneContent = rightPaneContent + "\n" + lipgloss.NewStyle().Foreground(green).Render(string(pretty))
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			userinputPaneStyle.Render(fmt.Sprintf("%4s", m.dockeruserinput.View())),
			leftPaneStyle.Render(""),
		),
		rightPaneStyle.Render(rightPaneContent),
	)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
