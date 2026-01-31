package state

import (
	"fmt"
	"slices"
)

// pieceTypes - valid piece types
var pieceTypes = []string{
	"rook",
	"bishop",
	"knight",
	"king",
	"queen",
	"pawn",
}

// piecesWorth - worth of each piece type
var piecesWorth = map[string]int{
	"rook":   5,
	"bishop": 3,
	"knight": 3,
	"queen":  9,
	"pawn":   1,
}

// piecesSymbols - unicode symbols for each piece type and color
var piecesSymbols = map[string](map[string]string){
	"rook":   map[string]string{"white": "\u2656", "black": "\u265C"},
	"bishop": map[string]string{"white": "\u2657", "black": "\u265D"},
	"king":   map[string]string{"white": "\u2654", "black": "\u265A"},
	"knight": map[string]string{"white": "\u2658", "black": "\u265E"},
	"queen":  map[string]string{"white": "\u2655", "black": "\u265B"},
	"pawn":   map[string]string{"white": "\u2664", "black": "\u2660"},
}

// Piece - represents a chess piece
// didn't use an interface because it didn't give much benefit (in go interfaces are implicit)
type Piece struct {
	Color              string
	Type               string
	Pos                Position
	HasMoved           bool
	Worth              int
	possibleMovesCache []*Move
}

// direction - helper struct for ray and step moves
type direction struct {
	Dx int
	Dy int
}

// typeIsValid - checks whether the piece type is valid (rook, bishop, knight, king, queen, pawn)
func typeIsValid(typ string) bool {
	return slices.Contains(pieceTypes, typ)
}

// colorIsValid - checks whether the piece color is valid (white or black)
func colorIsValid(color string) bool {
	return color == "white" || color == "black"
}

// createPiece - creates a new piece if the type and color are valid, otherwise returns nil
func createPiece(color string, typ string, pos Position) (*Piece, error) {
	if typeIsValid(typ) && colorIsValid(color) {
		return &Piece{
			Color:              color,
			Type:               typ,
			Pos:                pos,
			HasMoved:           false,
			Worth:              piecesWorth[typ],
			possibleMovesCache: []*Move{},
		}, nil
	}
	return nil, fmt.Errorf("invalid piece")
}

// GetPossibleMoves - get all *possible* moves of the pieces on the given board
// switches over piece type and calls the appropriate helper function
func (p *Piece) GetPossibleMoves(board *Board) ([]*Move, error) {
	var err error
	if len(p.possibleMovesCache) != 0 {
		return p.possibleMovesCache, nil
	}
	switch p.Type {
	case "rook":
		p.possibleMovesCache, err = p.getPossibleMovesRook(board)

	case "bishop":
		p.possibleMovesCache, err = p.getPossibleMovesBishop(board)

	case "knight":
		p.possibleMovesCache, err = p.getPossibleMovesKnight(board)

	case "queen":
		p.possibleMovesCache, err = p.getPossibleMovesQueen(board)

	case "king":
		p.possibleMovesCache, err = p.getPossibleMovesKing(board)

	case "pawn":
		p.possibleMovesCache, err = p.getPossibleMovesPawn(board)

	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p.possibleMovesCache, nil
}

// rayMoves - helper function for pieces that move in rays (rook, bishop, queen)
func (p *Piece) rayMoves(board *Board, directions *[]direction) ([]*Move, error) {
	possibleMoves := []*Move{}
	for _, d := range *directions {
		x, y := p.Pos.X+d.Dx, p.Pos.Y+d.Dy
		dst := Position{X: x, Y: y}
		for board.isInBounds(&dst) {
			target, err := board.GetPiece(&dst)
			if err != nil {
				return nil, err
			}
			if target == nil {
				possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
			} else {
				if target.Color != p.Color {
					possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
				}
				break
			}
			dst = Position{X: dst.X + d.Dx, Y: dst.Y + d.Dy}
		}
	}
	return possibleMoves, nil
}

// stepMoves - helper function for pieces that move in steps (king, knight)
func (p *Piece) stepMoves(board *Board, directions *[]direction) ([]*Move, error) {
	possibleMoves := []*Move{}
	for _, d := range *directions {
		x, y := p.Pos.X+d.Dx, p.Pos.Y+d.Dy
		dst := Position{X: x, Y: y}
		if board.isInBounds(&dst) {
			target, err := board.GetPiece(&dst)
			if err != nil {
				return nil, err
			}
			if target == nil {
				possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
			} else {
				if target.Color != p.Color {
					possibleMoves = append(possibleMoves, CreateMove(p.Pos, dst))
				}
			}
		}

	}
	return possibleMoves, nil
}

// clearCache - clears the possible moves cache
func (p *Piece) clearCache() {
	p.possibleMovesCache = []*Move{}
}

// moveTo - moves the piece to the given position and marks it as having moved
func (p *Piece) moveTo(pos Position) {
	p.Pos = pos
	p.HasMoved = true
	p.clearCache()
}

func (p *Piece) String() string {
	return fmt.Sprintf("{color: %v, type: %v, pos: %v, hasMoved: %v}", p.Color, p.Type, p.Pos, p.HasMoved)
}

// copy - creates a opy of the piece  used for simulating moves
func (p *Piece) copy() *Piece {
	if p == nil {
		fmt.Println("piece is nil")
	}
	return &Piece{
		Color:    p.Color,
		Type:     p.Type,
		Pos:      p.Pos,
		HasMoved: p.HasMoved,
		Worth:    p.Worth,
	}
}

func (p *Piece) Symbol() string {
	return piecesSymbols[p.Type][p.Color]
}
