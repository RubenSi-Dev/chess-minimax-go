package main

import (
	//"fmt"
	"fmt"

	color "github.com/fatih/color"
	chess "github.com/spunker/chess/state"
)

var whiteSquare = color.BgRGB(210, 210, 180)
var blueSquare = color.New(color.BgBlue)
var greenSquare = color.New(color.BgGreen)
var selectedSquare = color.New(color.BgRed)
var cursorWhiteSquare = color.BgRGB(245, 245, 245)
var cursorBlueSquare = color.New(color.BgCyan)
var cursorSelectedSquare = color.New(color.BgHiRed)

var spacingBefore = "            "

func printSquareColored(square *chess.Piece, col *color.Color) (result string) {
	result = ""
	if square != nil {
		return col.Sprintf(color.BlackString("   %v   "), square.Symbol())
	} else {
		return col.Sprint("       ")
	}
}

func printLine(rank []*chess.Piece, white bool, number int, selected int, cursored int) (result string) {
	if number == -1 {
		result = greenSquare.Sprint("     ")
	} else {
		result = greenSquare.Sprintf("  %v  ", number)
	}

	var color *color.Color
	white = !white
	for i, piece := range rank {
		white = !white
		if white && selected != i && cursored != i {
			color = whiteSquare
		} else if !white && selected != i && cursored != i {
			color = blueSquare
		} else if selected == i && cursored != i {
			color = selectedSquare
		} else if selected != i && cursored == i {
			if white {
				color = cursorWhiteSquare
			} else {
				color = cursorBlueSquare
			}
		} else if selected == i && cursored == i {
			color = cursorSelectedSquare
		}

		//fmt.Printf("color used: ", color)
		result += printSquareColored(piece, color)
	}

	if number == -1 {
		result += greenSquare.Sprint("     ")
	} else {
		result += greenSquare.Sprintf("  %v  ", number)
	}
	return
}

func printRank(rank []*chess.Piece, white bool, number int, sel []chess.Position, cursor chess.Position, statistic string) (result string) {
	var selected int
	if len(sel) == 1 && sel[0].Y == number-1 {
		selected = sel[0].X
	} else if len(sel) == 2 {
		if sel[0].Y == number-1 {
			selected = sel[0].X
		} else if sel[1].Y == number-1 {
			selected = sel[1].X
		} else {
			selected = -1
		}
	} else {
		selected = -1
	}

	var cursored int
	if cursor.Y == number-1 {
		cursored = cursor.X
	} else {
		cursored = -1
	}

	emptyrank := []*chess.Piece{nil, nil, nil, nil, nil, nil, nil, nil}

	//first line
	result += spacingBefore + printLine(emptyrank, white, -1, selected, cursored) + "\n"

	//second line
	result += spacingBefore + printLine(rank, white, number, selected, cursored) + fmt.Sprintln(statistic)

	//third line
	result += spacingBefore + printLine(emptyrank, white, -1, selected, cursored) + "\n"
	return
}

func printRankReverse(rank []*chess.Piece, white bool, number int, sel []chess.Position, cursor chess.Position, statistic string) (result string) {
	var selected int
	if len(sel) == 1 && sel[0].Y == number-1 {
		selected = sel[0].X
	} else if len(sel) == 2 {
		if sel[0].Y == number-1 {
			selected = sel[0].X
		} else if sel[1].Y == number-1 {
			selected = sel[1].X
		} else {
			selected = -1
		}
	} else {
		selected = -1
	}

	var cursored int
	if cursor.Y == number-1 {
		cursored = cursor.X
	} else {
		cursored = -1
	}

	emptyrank := []*chess.Piece{nil, nil, nil, nil, nil, nil, nil, nil}

	//first line
	result += spacingBefore + printLine(emptyrank, white, -1, selected, cursored) + "\n"

	//second line
	result += spacingBefore + printLine(rank, white, number, selected, cursored) + fmt.Sprintln(statistic)

	//third line
	result += spacingBefore + printLine(emptyrank, white, -1, selected, cursored) + "\n"
	return
}

func (m model) boardView() (result string) {
	if m.menu.playerColor == "black" {
		return m.boardViewBlack()
	}

	var lastMoveString string
	if len(m.game.State.PreviousMoves) > 0 {
		lastMoveString = m.game.State.PreviousMoves[len(m.game.State.PreviousMoves)-1].ToAlgebraic()
	} else {
		lastMoveString = ""
	}

	result += "\n"
	result += spacingBefore + greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += spacingBefore + greenSquare.Sprintln("                                                                  ")
	result += printRank(m.game.State.Board.Grid[7], false, 8, m.selected, m.cursor, fmt.Sprintf("       advantage for white: %v", GetMaterialStats(m.game.State.Board).GetAdvantage("white")))
	result += printRank(m.game.State.Board.Grid[6], true, 7, m.selected, m.cursor, fmt.Sprintf("       bot evaluation:      %v", botEvaln))
	result += printRank(m.game.State.Board.Grid[5], false, 6, m.selected, m.cursor, fmt.Sprintf("       to move:             %v", m.game.State.Turn))
	result += printRank(m.game.State.Board.Grid[4], true, 5, m.selected, m.cursor, fmt.Sprintf("       last move:           %v", lastMoveString))
	result += printRank(m.game.State.Board.Grid[3], false, 4, m.selected, m.cursor, "")
	result += printRank(m.game.State.Board.Grid[2], true, 3, m.selected, m.cursor, "")
	result += printRank(m.game.State.Board.Grid[1], false, 2, m.selected, m.cursor, "")
	result += printRank(m.game.State.Board.Grid[0], true, 1, m.selected, m.cursor, "")
	result += spacingBefore + greenSquare.Sprintln("                                                                  ")
	result += spacingBefore + greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += "\n"
	return result
}

func (m model) boardViewBlack() (result string) {
	result += "\n"
	result += spacingBefore + greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += spacingBefore + greenSquare.Sprintln("                                                                  ")
	result += printRankReverse(m.game.State.Board.Grid[0], false, 1, m.selected, m.cursor, fmt.Sprintf("       advantage for white: %v", GetMaterialStats(m.game.State.Board).GetAdvantage("white")))
	result += printRankReverse(m.game.State.Board.Grid[1], true, 2, m.selected, m.cursor, fmt.Sprintf("       bot evaluation:      %v", botEvaln))
	result += printRankReverse(m.game.State.Board.Grid[2], false, 3, m.selected, m.cursor, "")
	result += printRankReverse(m.game.State.Board.Grid[3], true, 4, m.selected, m.cursor, "")
	result += printRankReverse(m.game.State.Board.Grid[4], false, 5, m.selected, m.cursor, "")
	result += printRankReverse(m.game.State.Board.Grid[5], true, 6, m.selected, m.cursor, "")
	result += printRankReverse(m.game.State.Board.Grid[6], false, 7, m.selected, m.cursor, "")
	result += printRankReverse(m.game.State.Board.Grid[7], true, 8, m.selected, m.cursor, "")
	result += spacingBefore + greenSquare.Sprintln("                                                                  ")
	result += spacingBefore + greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += "\n"
	return result
}
