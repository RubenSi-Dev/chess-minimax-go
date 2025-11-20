package state

import (
	"fmt"
	"slices"
)

var pieceTypes = []string{
	"rook",
	"bishop",
	"knight",
	"king",
	"queen",
	"pawn",
}

var piecesWorth = map[string]int{
	"rook": 5,
	"bishop": 3,
	"knight": 3,
	"queen": 9,
	"pawn": 1,
}

var piecesSymbols = map[string](map[string]string) {
	"rook": map[string]string{"white": "\u2656", "black": "\u265C"},
	"bishop": map[string]string{"white": "\u2657", "black": "\u265D"},
	"king": map[string]string{"white": "\u2654", "black": "\u265A"},
	"knight": map[string]string{"white": "\u2658", "black": "\u265E"},
	"queen": map[string]string{"white": "\u2655", "black": "\u265B"},
	"pawn": map[string]string{"white": "\u2664", "black": "\u2660"},
}

type Piece struct {
	Color string
	Type string
	Pos Position
	HasMoved bool
	Worth int
	possibleMovesCache []*Move
}

type direction struct {
	Dx int
	Dy int
}

func typeIsValid(typ string) bool {
	return slices.Contains(pieceTypes, typ)
}

func colorIsValid(color string) bool {
	return color == "white" || color == "black"
}

func createPiece(color string, typ string, pos Position) *Piece {
	if typeIsValid(typ) && colorIsValid(color) {
		return &Piece{
			Color: color,
			Type: typ,
			Pos: pos,
			HasMoved: false,
			Worth: piecesWorth[typ],
			possibleMovesCache: []*Move{},
		}
	}
	return nil
}

func (p *Piece) GetPossibleMoves(board *Board) []*Move {
	if len(p.possibleMovesCache) != 0 {
		return p.possibleMovesCache
	}
	switch p.Type {
	case "rook":
		p.possibleMovesCache = p.getPossibleMovesRook(board)
	case "bishop":
		p.possibleMovesCache = p.getPossibleMovesBishop(board)
	case "knight":
		p.possibleMovesCache = p.getPossibleMovesKnight(board)
	case "queen":
		p.possibleMovesCache = p.getPossibleMovesQueen(board)
	case "king":
		p.possibleMovesCache = p.getPossibleMovesKing(board)
	case "pawn":
		p.possibleMovesCache = p.getPossibleMovesPawn(board)
	default:
		return nil
	}
	return p.possibleMovesCache
}

func (p *Piece) rayMoves(board *Board, directions *[]direction) (possibleMoves []*Move) {
	possibleMoves = []*Move{}
	for _, d := range *directions {
		x, y := p.Pos.X + d.Dx, p.Pos.Y + d.Dy
		dst := Position{X: x, Y: y}
		for board.isInBounds(&dst) {
			target := board.GetPiece(&dst)
			if target == nil {
				possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
			} else {
				if (target.Color != p.Color) {
					possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
				}
				break
			}
		dst = Position{X: dst.X+d.Dx, Y: dst.Y+d.Dy}
		}
	}
	return
}

func (p *Piece) stepMoves(board *Board, directions *[]direction) (possibleMoves []*Move) {
	possibleMoves = []*Move{}
	for _, d := range *directions {
		x, y := p.Pos.X + d.Dx, p.Pos.Y + d.Dy
		dst := Position{X: x, Y: y}
		if board.isInBounds(&dst) {
			target := board.GetPiece(&dst)
			if (target == nil) {
				possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
			} else {
				if (target.Color != p.Color) {
					possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
				}
			}
		}

	}
	return
}

func (p *Piece) clearCache() {
	p.possibleMovesCache = []*Move{}
}

func (p *Piece) moveTo(pos Position) {
	p.Pos = pos
	p.HasMoved = true
	p.clearCache()
}

func (p *Piece) String() string {
	return fmt.Sprintf("{color: %v, type: %v, pos: %v, hasMoved: %v}", p.Color, p.Type, p.Pos, p.HasMoved)
}

func (p *Piece) copy() *Piece {
	if p == nil {
		fmt.Println("piece is nil")
	}
	return &Piece{
		Color: p.Color,
		Type: p.Type,
		Pos: p.Pos,
		HasMoved: p.HasMoved,
		Worth: p.Worth,
	}
}

func (p *Piece) Symbol() string {
	return piecesSymbols[p.Type][p.Color]
}
