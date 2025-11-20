package state

import (
	//"fmt"
	"reflect"
)

var setups = []string{
	"default",
	"promotion",
	"clear",
	"castling",
}

var algebraicLetters = []string{"A", "B", "C", "D", "E", "F", "G", "H"}

type grid [][]*Piece

type Board struct {
	Grid grid
	piecesCache []*Piece
	squaresControlledCache map[string]([]*Position)
	kingsPositionCache map[string]*Position
}

func createBoard(setup string) (result *Board) {
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
	}
	result.initBoard(setup)
	return
}

func (b *Board) clearCache() {
	for _, piece := range b.piecesCache {
		piece.clearCache()
	}
	b.piecesCache = []*Piece{}
	for k := range b.squaresControlledCache {
		b.squaresControlledCache[k]	= []*Position{}
	}
	for k := range b.kingsPositionCache {
		b.kingsPositionCache[k] = nil
	}
}

func (b *Board) initBoard(setup string) {
	b.Grid = make([][]*Piece, 8)
	for i := range b.Grid {
		b.Grid[i] = make([]*Piece, 8)
	}
	switch setup {
	case "clear":
		return
	case "default":
		b.defaultSetup()
	case "promotion":
		b.promotionSetup()
	case "castling":
		b.castlingSetup()
	}
}

func (b *Board) castlingSetup() {
	//kings
	b.placeNew("white", "king", Position{X: 4, Y: 0})
	b.placeNew("black", "king", Position{X: 4, Y: 7})

	b.placeNew("white", "rook", Position{X: 0, Y: 0})
	b.placeNew("white", "rook", Position{X: 7, Y: 0})
	b.placeNew("black", "rook", Position{X: 0, Y: 7})
	b.placeNew("black", "rook", Position{X: 7, Y: 7})
}

func (b *Board) promotionSetup() {
	//kings
	b.placeNew("white", "king", Position{X: 4, Y: 5})
	b.placeNew("black", "king", Position{X: 7, Y: 7})

	//pawns
	b.placeNew("white", "pawn", Position{X: 3, Y: 6})
	b.placeNew("black", "pawn", Position{X: 7, Y: 6})
}

func (b *Board) defaultSetup() {
	//rooks
	b.placeNew("white", "rook", Position{X: 0, Y: 0})
	b.placeNew("white", "rook", Position{X: 7, Y: 0})
	b.placeNew("black", "rook", Position{X: 0, Y: 7})
	b.placeNew("black", "rook", Position{X: 7, Y: 7})

	//knights
	b.placeNew("white", "knight", Position{X: 1, Y: 0})
	b.placeNew("white", "knight", Position{X: 6, Y: 0})
	b.placeNew("black", "knight", Position{X: 1, Y: 7})
	b.placeNew("black", "knight", Position{X: 6, Y: 7})

	//bishop
	b.placeNew("white", "bishop", Position{X: 2, Y: 0})
	b.placeNew("white", "bishop", Position{X: 5, Y: 0})
	b.placeNew("black", "bishop", Position{X: 2, Y: 7})
	b.placeNew("black", "bishop", Position{X: 5, Y: 7})

	//queens
	b.placeNew("white", "queen", Position{X: 3, Y: 0})
	b.placeNew("black", "queen", Position{X: 3, Y: 7})

	//kings
	b.placeNew("white", "king", Position{X: 4, Y: 0})
	b.placeNew("black", "king", Position{X: 4, Y: 7})

	//pawns
	for i := range 8 {
		b.placeNew("white", "pawn", Position{X: i, Y: 1})
		b.placeNew("black", "pawn", Position{X: i, Y: 6})
	}
}


func (b *Board) GetPieces() ([]*Piece) {
	if len(b.piecesCache) == 0 { 
		b.piecesCache = []*Piece{}
		for _, rank := range b.Grid {
			for _, square := range rank{
				if (square != nil) {
					b.piecesCache = append(b.piecesCache, square);
				}
			}
		}
	}
	return b.piecesCache
}

func (b *Board) GetPiece(pos *Position) *Piece{
	if !b.isInBounds(pos) { return nil }
	return b.Grid[pos.Y][pos.X]
}

func (b *Board) FindPiece(typ string, color string) (result []*Position) {
	if (typ == "king" && b.kingsPositionCache[color] != nil) {
		return []*Position{b.kingsPositionCache[color]}
	} 
	result = []*Position{}
	for _, piece := range b.GetPieces() {
		if (piece.Color == color && piece.Type == typ) {
			result = append(result, &piece.Pos)
		}
	}
	if (typ == "king") {
		b.kingsPositionCache[color] = result[0]
	}
	return
}

func (b *Board) isInBounds(pos *Position) bool {
	return (pos.X < len(b.Grid) && pos.X >= 0) && (pos.Y < len(b.Grid[0]) && pos.Y >= 0)
}

func (b *Board) PlaceOn(piece *Piece, pos *Position) bool {
	if (!b.isInBounds(pos)) { return false }
	b.Grid[pos.Y][pos.X] = piece
	return true
}

func (b *Board) placeNew(color string, typ string, pos Position) (piece *Piece) {
	if (!b.isInBounds(&pos)) { return nil } 

	piece = createPiece(color, typ, pos)
	b.Grid[pos.Y][pos.X] = piece
	return
}

func (b *Board) RemoveFrom(pos *Position) {
	if (!b.isInBounds(pos)) { return }
	b.Grid[pos.Y][pos.X] = nil
}

func (b *Board) squaresControlledBy(color string) []*Position {
	if len(b.squaresControlledCache[color]) == 0 {
		for _,piece := range b.GetPieces() {
			if piece.Color == color {
				moves := piece.GetPossibleMoves(b)
				for _,move := range moves {
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

func (b *Board) copy() (copy *Board) {
	copy = createBoard("clear")
	for y, rank := range b.Grid {
		for x, square := range rank {
			if (square != nil) {
				copy.Grid[y][x] = b.Grid[y][x].copy()
			}
		}
	}
	return
}
