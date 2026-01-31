package state

import (
	"github.com/fatih/color"
)

var whiteSquare = color.New(color.BgWhite)
var blueSquare = color.New(color.BgBlue)
var greenSquare = color.New(color.BgGreen)

func printFullLine(rank []*Piece, white bool, number int) (result string) {
	result = greenSquare.Sprintf("  %v  ", number)
	if white {
		for i, square := range rank {
			result += printSquare(square, i%2 == 0)
		}
	} else {
		for i, square := range rank {
			result += printSquare(square, i%2 == 1)
		}
	}
	result += greenSquare.Sprintf("  %v  ", number)
	result += "\n"
	return
}

func printEmptyLine(white bool) (result string) {
	result = greenSquare.Sprint("     ")
	if white {
		for i := range 8 {
			result += printSquare(nil, i%2 == 0)
		}
	} else {
		for i := range 8 {
			result += printSquare(nil, i%2 == 1)
		}
	}
	result += greenSquare.Sprint("     ")
	result += "\n"
	return
}

func printRank(rank []*Piece, white bool, number int) (result string) {
	//first line
	result += printEmptyLine(white)

	//second line
	result += printFullLine(rank, white, number)

	//third line
	result += printEmptyLine(white)
	return
}

func printSquare(square *Piece, white bool) (result string) {
	result = ""
	if white {
		if square != nil {
			return whiteSquare.Sprintf(color.BlackString("   %v   "), square.Symbol())
		} else {
			return whiteSquare.Sprint("       ")
		}
	} else {
		if square != nil {
			return blueSquare.Sprintf(color.BlackString("   %v   "), square.Symbol())
		} else {
			return blueSquare.Sprint("       ")
		}
	}
}

func (b *Board) String() (result string) {
	result += greenSquare.Sprintln("                                                                  ")
	result += greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += greenSquare.Sprintln("                                                                  ")
	for k := range b.Grid {
		i := len(b.Grid) - k - 1
		rank := b.Grid[i]
		result += printRank(rank, i%2 == 0, i+1)
	}
	result += greenSquare.Sprintln("                                                                  ")
	result += greenSquare.Sprintln("        A      B      C      D      E      F      G      H        ")
	result += greenSquare.Sprintln("                                                                  ")
	return result
}
