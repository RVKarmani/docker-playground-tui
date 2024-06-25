package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

func initialModel() model {
	m := model{
		err:               nil,
		dockerhubresponse: DockerHubResponse{},
		websocketConn:     nil,
		websocketmessages: []WebsocketReqResp{},
	}

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
		m.err = msg
		return m, nil

	case DockerHubResponse:
		m.err = nil
		m.dockerhubresponse = msg
		m.dockeruserinput.Reset()
		m.dockeruserinput.Placeholder = "Enter command to send via websocket"
		return m, nil

	case WebsocketReqResp:
		m.websocketmessages = append(m.websocketmessages, msg)
		m.dockeruserinput.Reset()
		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			if (m.dockerhubresponse == DockerHubResponse{}) {
				// Means the user is not yet authenticated, go through that
				return m, getDockerHubDetails(m.dockeruserinput.Value())
			} else {
				// user is authenticated, go through command flow\
				return m, m.sendAndReceiveWebsocketMsg(m.dockeruserinput.Value())
			}
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

	leftPaneContent := ""

	if len(m.websocketmessages) > 0 {
		for _, wsMsgPair := range m.websocketmessages {
			leftPaneContent += "\n" + lipgloss.NewStyle().Foreground(hotPink).Render(wsMsgPair.inputstr) + "\n" + lipgloss.NewStyle().Foreground(green).Render(wsMsgPair.outputstr)
		}
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Center,
			userinputPaneStyle.Render(fmt.Sprintf("%4s", m.dockeruserinput.View())),
			leftPaneStyle.Render(leftPaneContent),
		),
		rightPaneStyle.Render(rightPaneContent),
	)
}

func main() {
	// Setup websocket connection here
	connection, _, err := websocket.DefaultDialer.Dial(websocketEndpoint, nil)
	if err != nil {
		fmt.Print("error")
	}

	connection.ReadMessage()

	initModel := initialModel()
	initModel.websocketConn = connection

	defer connection.Close()

	p := tea.NewProgram(initModel)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
