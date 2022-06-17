package full

import (
	"strconv"
	"strings"

	"github.com/smmr-software/mabel/internal/styles"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"

	"github.com/anacrolix/torrent"
)

type portStartupFailure struct {
	width, height int
	input         textinput.Model
	main          *model
}

func (m portStartupFailure) Init() tea.Cmd {
	return tick()
}

func (m portStartupFailure) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width - styles.BorderWindow.GetHorizontalBorderSize()
		m.height = msg.Height - styles.BorderWindow.GetHorizontalBorderSize()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "backspace":
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		default:
			prt, err := strconv.Atoi(m.input.Value())
			if err != nil {
				return m, reportError(err)
			}
			port := uint(prt)

			config := genMabelConfig(&port, m.main.logging)
			client, err := torrent.NewClient(config)
			if err != nil {
				return m, reportError(err)
			}

			m.main.client = client
			m.main.clientConfig = config
			m.main.width = m.width
			m.main.height = m.height

			return m.main, nil
		}
	case tickMsg:
		return m, tick()
	default:
		return m, nil
	}
}

func (m portStartupFailure) View() string {
	fullscreen := gloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Inherit(styles.BorderWindow)

	var body strings.Builder
	body.WriteString(styles.Bold.Render("Port Binding Failure"))
	body.WriteString("\nplease provide an unused port number for the client to bind with\n\n")
	body.WriteString(styles.BorderWindow.Render(m.input.View()))

	return fullscreen.Render(
		gloss.Place(
			m.width, m.height,
			gloss.Center, gloss.Center,
			body.String(),
		),
	)
}

func initialPortStartupFailure(parent *model) portStartupFailure {
	input := textinput.New()
	input.Width = 32
	input.Focus()

	return portStartupFailure{input: input, main: parent}
}
