package main 

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	//color "github.com/fatih/color"
	chess "github.com/spunker/chess/state"
	ai "github.com/spunker/chess/ai"
)

var botEvaln float64

type model struct {
	game *Game
	cursor chess.Position
	selected []chess.Position
}


func initialModel() model {
	return model{
		game: StartGame("default"),
		selected: []chess.Position{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type BotMoveMsg struct {
	move *chess.Move
	score float64
}

func getBotMove(s *chess.State, depth int) tea.Cmd {
	return func() tea.Msg {
		move, score := ai.SelectMove(s, depth)
		return BotMoveMsg{
			move: move,
			score: score,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	//key press
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.cursor.Y < 7 {
				m.cursor.Y++
			}

		case "down", "j":
			if m.cursor.Y > 0 {
				m.cursor.Y--
			}

		case "left", "h":
			if m.cursor.X > 0 {
				m.cursor.X--
			}

		case "right", "l":
			if m.cursor.X < 7 {
				m.cursor.X++
			}

		case "esc":
			m.selected = []chess.Position{}

		case "enter", " ":
			m.selected = append(m.selected, m.cursor)
			if len(m.selected) >= 2 {
				ok := m.game.PlayMove(&chess.Move{
					From: m.selected[0],
					To: m.selected[1],
				})
				m.selected = []chess.Position{}
				if ok {
					return m, getBotMove(m.game.State, 3)
				}
			}
		}

	case BotMoveMsg:
		m.game.PlayMove(msg.move)
		botEvaln = msg.score
		return m, nil
	}

	return m, nil
}


func (m model) View() string {
	return BoardString(m.game.State.Board, m.selected, m.cursor)
}

func StartTui() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
