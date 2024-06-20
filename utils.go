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
)

const (
	hotPink = lipgloss.Color("#FF06B7")
	green   = lipgloss.Color("#077B8A")
)

var (
	userinputPaneStyle = lipgloss.NewStyle().
				Width(60).
				Height(5).
				Align(lipgloss.Center, lipgloss.Left).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))

	leftPaneStyle = lipgloss.NewStyle().
			Width(60).
			Height(25).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

	rightPaneStyle = lipgloss.NewStyle().
			Width(60).
			Height(30).
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

type errMsg struct {
	err error
}

func (e errMsg) Error() string { return e.err.Error() }
