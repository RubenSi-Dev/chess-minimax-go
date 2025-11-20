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
	
	var color *color.Color;
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
	var selected int;
	if len(sel) == 1 && sel[0].Y == number - 1 {
		selected = sel[0].X
	} else if len(sel) == 2 {
		if sel[0].Y == number - 1 {
			selected = sel[0].X
		} else if sel[1].Y == number - 1 {
			selected = sel[1].X
		} else {
			selected = -1
		}
	} else {
		selected = -1
	}

	var cursored int;
	if cursor.Y == number - 1 {
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

func BoardString(b *chess.Board, sel []chess.Position, cursor chess.Position) (result string) {
	result += "\n"
	result += spacingBefore + greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += spacingBefore + greenSquare.Sprintln("                                                                  ")
	result += printRank(b.Grid[7], false, 8, sel, cursor, fmt.Sprintf("       advantage for white: %v", GetMaterialStats(b).GetAdvantage("white")))
	result += printRank(b.Grid[6], true, 7, sel, cursor, fmt.Sprintf("       bot evaluation:      %v", botEvaln))
	result += printRank(b.Grid[5], false, 6, sel, cursor, "")
	result += printRank(b.Grid[4], true, 5, sel, cursor, "")
	result += printRank(b.Grid[3], false, 4, sel, cursor, "")
	result += printRank(b.Grid[2], true, 3, sel, cursor, "")
	result += printRank(b.Grid[1], false, 2, sel, cursor, "")
	result += printRank(b.Grid[0], true, 1, sel, cursor, "")
	result += spacingBefore + greenSquare.Sprintln("                                                                  ")
	result += spacingBefore + greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += "\n"
	return result
}
