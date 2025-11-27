package main

import (
	"fmt"
	"math"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	color "github.com/fatih/color"

	//color "github.com/fatih/color"
	ai "github.com/spunker/chess/ai"
	chess "github.com/spunker/chess/state"
)

var botEvaln float64

type Menu struct {
	playerColor string
	setup string
	botDepth int
	weights ai.Weights
}

type model struct {
	game        *Game
	cursor      chess.Position
	selected    []chess.Position
	inMenu      bool
	menu Menu	
	menuCursor int
}

func initialModel() model {
	return model{
		inMenu:      true,
		menu: Menu{
			playerColor: "white",
			setup: "default",
			botDepth: 3,
			weights: ai.Weights{
				Material: 2.0,
				Mobility: 0.5,
			},
		},
		menuCursor:  0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type BotMoveMsg struct {
	move  *chess.Move
	score float64
}

func (m model) getBotMove(s *chess.State, depth int) tea.Cmd {
	return func() tea.Msg {
		move, score := ai.SelectMove(s, depth, &m.menu.weights)
		return BotMoveMsg{
			move:  move,
			score: score,
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inMenu {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "left", "h":
				switch m.menuCursor {
				case 0: // playercolor
					if m.game == nil {
						if m.menu.playerColor == "white" {
							m.menu.playerColor = "black"
						} else {
							m.menu.playerColor = "white"
						}

					}

				case 1: // setup
					choices := []string{"default", "castling", "promotion", "clear"}
					if m.game == nil {
						for i, setup := range choices  {
							if setup == m.menu.setup {
								if i == 0 { 
									m.menu.setup = choices[len(choices)-1]
								} else {
									m.menu.setup = choices[i-1]
								}
							}
						}
					}

				case 2: // botDepth
					if m.menu.botDepth > 0 {
						m.menu.botDepth--
					}

				case 3: // weight material
					m.menu.weights.Material -= 0.1
					m.menu.weights.Material = math.Round(m.menu.weights.Material * 10)/10

				case 4: // weight mobility
					m.menu.weights.Mobility -= 0.1
					m.menu.weights.Mobility = math.Round(m.menu.weights.Mobility * 10)/10
				}

			case "right", "l":
				switch m.menuCursor {
				case 0: // playercolor
					if m.menu.playerColor == "white" {
						m.menu.playerColor = "black"
					} else {
						m.menu.playerColor = "white"
					}

				case 1: // setup
					choices := []string{"default", "castling", "promotion", "clear"}
					if m.game == nil {
						for i, setup := range choices  {
							if setup == m.menu.setup {
								if i == len(choices) - 1 {
									m.menu.setup = choices[0]
								} else {
									m.menu.setup = choices[i+1]
								}
							}
						}
					}

				case 2: // botDepth
					if m.menu.botDepth < 5 {
						m.menu.botDepth++
					}

				case 3: // weight material
					m.menu.weights.Material += 0.1
					m.menu.weights.Material = math.Round(m.menu.weights.Material * 10)/10

				case 4: // weight mobility
					m.menu.weights.Mobility += 0.1
					m.menu.weights.Mobility = math.Round(m.menu.weights.Mobility * 10)/10
				}
				
			case "up", "k":
				if m.menuCursor > 0 {
					m.menuCursor--
				}

			case "down", "j":
				if m.menuCursor < 4 {
					m.menuCursor++
				}

			case "enter", " ":
				m.inMenu = false
				if m.game == nil {
					m.game = StartGame(m.menu.setup)
					m.cursor = chess.Position{
						X: 4,
						Y: 4,
					}
				}
				if m.menu.playerColor == "black" {
					return m, m.getBotMove(m.game.State, m.menu.botDepth)
				}
			}
		}

	} else {
		switch msg := msg.(type) {
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

			case "i":
				m.inMenu = true

			case "enter", " ":
				if (m.game.State.Turn == m.menu.playerColor) {
					m.selected = append(m.selected, m.cursor)
					if len(m.selected) >= 2 {
						ok := m.game.PlayMove(&chess.Move{
							From: m.selected[0],
							To:   m.selected[1],
						})
						m.selected = []chess.Position{}
						if ok {
							return m, m.getBotMove(m.game.State, m.menu.botDepth)
						}
					}

				}
			}

		case BotMoveMsg:
			m.game.PlayMove(msg.move)
			botEvaln = msg.score
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.inMenu {
		return m.menuView()
	}
	return m.boardView()
}

func (m model) getCursorString() (result map[string]string) {
	result = map[string]string{
		"playerColor": "  ",
		"setup": "  ",
		"botDepth": "  ",
		"materialWeights": "  ",
		"mobilityWeights": "  ",
	}
	switch m.menuCursor {
	case 0: 
		result["playerColor"] = " >"
	case 1:
		result["setup"] = " >"
	case 2: 
		result["botDepth"] = " >"
	case 3: 
		result["materialWeights"] = " >"
	case 4: 
		result["mobilityWeights"] = " >"
	}
	return
}

func (m model) menuView() (result string) {
	result = "\n"

	cursorString := m.getCursorString()

	result += fmt.Sprintf("%v   You play as:         < %v >\n", cursorString["playerColor"], m.menu.playerColor)

	result += "\n"
	result += fmt.Sprintf("%v   Setup:               < %v >\n", cursorString["setup"], m.menu.setup)

	result += "\n"
	if m.menu.botDepth == 0 {
		result += fmt.Sprintf("%v   Engine depth:        < disabled >\n", cursorString["botDepth"])
	} else if m.menu.botDepth >= 5 {
		result += fmt.Sprintf("%v   Engine depth:        < %v > 	%v\n", cursorString["botDepth"], m.menu.botDepth, color.RedString("(not recommended)"))
	} else {
		result += fmt.Sprintf("%v   Engine depth:        < %v > \n", cursorString["botDepth"], m.menu.botDepth)
	}

	result += "\n"

	result += "   Weights for bot\n"
	result += "\n"

	result += fmt.Sprintf("%v   material:            < %v > \n", cursorString["materialWeights"], m.menu.weights.Material)
	result += "\n"
	result += fmt.Sprintf("%v   mobility:            < %v > \n", cursorString["mobilityWeights"], m.menu.weights.Mobility)
	return
}

func StartTui() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
