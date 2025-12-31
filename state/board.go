package state

import (
	"fmt"
	"reflect"
)

// go doesn't have an enum type, so we use a slice of strings to represent possible setups
var setups = []string{
	"default",
	"promotion",
	"clear",
	"castling",
}

// used for algebraic notation
var algebraicLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H"}

// grid - 8x8 grid of pointers to pieces
type grid [][]*Piece

// Board - represents the chess board and its pieces
// caching to avoid recomputation when nothing has changed
type Board struct {
	Grid                   grid
	piecesCache            []*Piece
	squaresControlledCache map[string]([]*Position)
	kingsPositionCache     map[string]*Position
	isCopy                 bool
}

// createBoard - creates a new board with the given setup
func createBoard(setup string) (result *Board, err error) {
	result = &Board{
		nil,
		[]*Piece{},
		map[string]([]*Position){
			"white": []*Position{},
			"black": []*Position{},
		},
		map[string]*Position{
			"white": nil,
			"black": nil,
		},
		false,
	}
	err = result.initBoard(setup)
	return
}

// clearCache - clears cached pieces and squares controlled
// usually called when the board changes (a move is played)
func (b *Board) clearCache() {
	for _, piece := range b.piecesCache {
		piece.clearCache()
	}
	b.piecesCache = []*Piece{}
	for k := range b.squaresControlledCache {
		b.squaresControlledCache[k] = []*Position{}
	}
}

func (b *Board) clearKingsPosCash(color string) {
	b.squaresControlledCache[color] = nil
}

// initBoard - initializes the board with the given setup, and all setups are defined below
func (b *Board) initBoard(setup string) error {
	b.Grid = make([][]*Piece, 8)
	for i := range b.Grid {
		b.Grid[i] = make([]*Piece, 8)
	}
	switch setup {
	case "clear":
		return nil
	case "default":
		return b.defaultSetup()
	case "promotion":
		return b.promotionSetup()
	case "castling":
		return b.castlingSetup()
	}
	return nil
}

func (b *Board) castlingSetup() error {
	//kings
	if _, err := b.placeNew("white", "king", Position{X: 4, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "king", Position{X: 4, Y: 7}); err != nil {
		return err
	}

	if _, err := b.placeNew("white", "rook", Position{X: 0, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("white", "rook", Position{X: 7, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "rook", Position{X: 0, Y: 7}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "rook", Position{X: 7, Y: 7}); err != nil {
		return err
	}
	return nil
}

func (b *Board) promotionSetup() error {
	//kings
	if _, err := b.placeNew("white", "king", Position{X: 4, Y: 5}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "king", Position{X: 7, Y: 7}); err != nil {
		return err
	}

	//pawns
	if _, err := b.placeNew("white", "pawn", Position{X: 3, Y: 6}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "pawn", Position{X: 7, Y: 6}); err != nil {
		return err
	}
	return nil
}

func (b *Board) defaultSetup() error {
	//rooks
	if _, err := b.placeNew("white", "rook", Position{X: 0, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("white", "rook", Position{X: 7, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "rook", Position{X: 0, Y: 7}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "rook", Position{X: 7, Y: 7}); err != nil {
		return err
	}

	//knights
	if _, err := b.placeNew("white", "knight", Position{X: 1, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("white", "knight", Position{X: 6, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "knight", Position{X: 1, Y: 7}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "knight", Position{X: 6, Y: 7}); err != nil {
		return err
	}

	//bishop
	if _, err := b.placeNew("white", "bishop", Position{X: 2, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("white", "bishop", Position{X: 5, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "bishop", Position{X: 2, Y: 7}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "bishop", Position{X: 5, Y: 7}); err != nil {
		return err
	}

	//queens
	if _, err := b.placeNew("white", "queen", Position{X: 3, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "queen", Position{X: 3, Y: 7}); err != nil {
		return err
	}

	//kings
	if _, err := b.placeNew("white", "king", Position{X: 4, Y: 0}); err != nil {
		return err
	}
	if _, err := b.placeNew("black", "king", Position{X: 4, Y: 7}); err != nil {
		return err
	}

	//pawns
	for i := range 8 {
		if _, err := b.placeNew("white", "pawn", Position{X: i, Y: 1}); err != nil {
			return err
		}
		if _, err := b.placeNew("black", "pawn", Position{X: i, Y: 6}); err != nil {
			return err
		}
	}
	return nil
}

// GetPieces - returns a list of all pieces on the board
func (b *Board) GetPieces() []*Piece {
	if len(b.piecesCache) == 0 {
		b.piecesCache = []*Piece{}
		for _, rank := range b.Grid {
			for _, square := range rank {
				if square != nil {
					b.piecesCache = append(b.piecesCache, square)
				}
			}
		}
	}
	return b.piecesCache
}

// GetPiece - returns the piece at the given position, or nil if there is no piece
func (b *Board) GetPiece(pos *Position) *Piece {
	if !b.isInBounds(pos) {
		return nil
	}
	return b.Grid[pos.Y][pos.X]
}

// FindPiece - returns a list of positions of pieces of the given type and color
func (b *Board) FindPiece(typ string, color string) (result []*Position) {
	if typ == "king" && b.kingsPositionCache[color] != nil {
		return []*Position{b.kingsPositionCache[color]}
	}
	result = []*Position{}
	for _, piece := range b.GetPieces() {
		if piece.Color == color && piece.Type == typ {
			result = append(result, &piece.Pos)
		}
	}
	if typ == "king" {
		b.kingsPositionCache[color] = result[0]
	}
	return
}

// isInBounds - checks whether the given position is within the bounds of the board
func (b *Board) isInBounds(pos *Position) bool {
	return (pos.X < len(b.Grid) && pos.X >= 0) && (pos.Y < len(b.Grid[0]) && pos.Y >= 0)
}

// PlaceOn - places the given piece on the given position
func (b *Board) PlaceOn(piece *Piece, pos *Position) bool {
	if !b.isInBounds(pos) {
		return false
	}
	b.Grid[pos.Y][pos.X] = piece
	return true
}

// placeNew - creates a new piece of the given type and color at the given position
func (b *Board) placeNew(color string, typ string, pos Position) (*Piece, error) {
	if !b.isInBounds(&pos) {
		return nil, fmt.Errorf("piece not in bounds")
	}

	piece, err := createPiece(color, typ, pos)
	if err != nil {
		return nil, fmt.Errorf("creation of piece failed")
	}
	b.Grid[pos.Y][pos.X] = piece
	return piece, nil
}

// RemoveFrom - removes the piece at the given position
func (b *Board) RemoveFrom(pos *Position) {
	if !b.isInBounds(pos) {
		return
	}
	b.Grid[pos.Y][pos.X] = nil
}

// squaresControlledBy - returns a list of positions controlled by the given color (the squares the color can move to)
func (b *Board) squaresControlledBy(color string) []*Position {
	if len(b.squaresControlledCache[color]) == 0 {
		for _, piece := range b.GetPieces() {
			if piece.Color == color {
				moves := piece.GetPossibleMoves(b)
				for _, move := range moves {
					b.squaresControlledCache[color] = append(b.squaresControlledCache[color], &move.To)
				}
			}
		}
	}
	return b.squaresControlledCache[color]
}

func (this *Board) Equal(other *Board) bool {
	return reflect.DeepEqual(this.Grid, other.Grid)
}

// copy - creates a deep copy of the board (also called by state.Copy())
func (b *Board) copy() (*Board, error) {
	copy, err := createBoard("clear")
	if err != nil {
		return nil, err
	}
	for y, rank := range b.Grid {
		for x, square := range rank {
			if square != nil {
				copy.Grid[y][x] = b.Grid[y][x].copy()
			}
		}
	}
	copy.isCopy = true
	return copy, nil
}
