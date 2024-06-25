package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

const (
	hotPink = lipgloss.Color("#FF06B7")
	green   = lipgloss.Color("#077B8A")
)

type WebsocketReqResp struct {
	inputstr  string
	outputstr string
}

const websocketEndpoint string = "wss://echo.websocket.org"

var (
	userinputPaneStyle = lipgloss.NewStyle().
				Width(60).
				Height(1).
				Align(lipgloss.Center, lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	leftPaneStyle = lipgloss.NewStyle().
			Width(60).
			Height(29).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

	rightPaneStyle = lipgloss.NewStyle().
			Width(60).
			Height(32).
			Align(lipgloss.Left, lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))
)

type DockerHubResponse struct {
	ID            string    `json:"id"`
	UUID          string    `json:"uuid"`
	Username      string    `json:"username"`
	FullName      string    `json:"full_name"`
	Location      string    `json:"location"`
	Company       string    `json:"company"`
	ProfileURL    string    `json:"profile_url"`
	DateJoined    time.Time `json:"date_joined"`
	GravatarURL   string    `json:"gravatar_url"`
	GravatarEmail string    `json:"gravatar_email"`
	Type          string    `json:"type"`
}

type model struct {
	dockeruserinput   textinput.Model
	err               error
	dockerhubresponse DockerHubResponse
	websocketConn     *websocket.Conn
	websocketmessages []WebsocketReqResp
}

func getDockerHubDetails(dockerUsername string) tea.Cmd {
	return func() tea.Msg {
		response, err := http.Get(fmt.Sprintf("https://hub.docker.com/v2/users/%s", dockerUsername))

		if err != nil {
			return errMsg{err}
		}

		if response.StatusCode != 200 {
			return errMsg{fmt.Errorf("user not found, check your username")}
		}

		responseData, err := io.ReadAll(response.Body)

		if err != nil {
			return errMsg{err}
		}

		var dockerHubResponse DockerHubResponse
		json.Unmarshal(responseData, &dockerHubResponse)

		return dockerHubResponse
	}
}

func (m model) sendAndReceiveWebsocketMsg(inputCommand string) tea.Cmd {
	return func() tea.Msg {
		err := m.websocketConn.WriteMessage(websocket.TextMessage, []byte(inputCommand))

		if err != nil {
			return errMsg{err}
		}

		_, message, err := m.websocketConn.ReadMessage()

		if err != nil {
			return errMsg{err}
		}

		return WebsocketReqResp{
			inputstr:  inputCommand,
			outputstr: string(message),
		}
	}
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }
