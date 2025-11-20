package state

import (
	"slices"
	//"fmt"
)

type State struct {
	Board *Board
	Turn string
	LastMove *Move
	possibleMovesCache []*Move
	legalMovesCache []*Move
	legalMovesOrderedCache []*Move
}

func CreateState(setup string) *State {
	if !(slices.Contains(setups, setup)) { setup = "clear" }
	return &State{
		Board: createBoard(setup),
		Turn: "white",
		LastMove: nil,
		possibleMovesCache: []*Move{},
		legalMovesCache: []*Move{},
		legalMovesOrderedCache: []*Move{},
	}	
}

func (s *State) clearCache() {
	s.possibleMovesCache = []*Move{}
	s.legalMovesCache = []*Move{}
	s.legalMovesOrderedCache = []*Move{}
	s.Board.clearCache()
}

func (s *State) switchTurn() string {
	white := s.Turn == "white"
	if (white) {
		s.Turn = "black"
	} else {
		s.Turn = "white"
	}
	return s.Turn
}

func (s *State) GetPossibleMoves() []*Move {
	if (len(s.possibleMovesCache) != 0) { return s.possibleMovesCache }

	pieces := s.Board.GetPieces() 
	for _, piece := range pieces {
		if piece.Color == s.Turn {
			moves := piece.GetPossibleMoves(s.Board)
			//fmt.Printf("calculated moves for %v: %v\n", piece.Symbol(), moves)
			s.possibleMovesCache = append(s.possibleMovesCache, moves...)
		}
	}
	return s.possibleMovesCache
}


func (this *State) applyMoveBool(move *Move) bool {
	piece := this.Board.GetPiece(&move.From) 
	if piece != nil {
		if move.Promotion == "queen" {
			this.Board.RemoveFrom(&move.From)
			newPiece := this.Board.placeNew(piece.Color, move.Promotion, move.To)
			newPiece.HasMoved = true
			this.LastMove = move
			//fmt.Println(this.Board)
			return true
		}	else {
			if piece.Type == "king" && !piece.HasMoved {
				if move.From.X - move.To.X == -2 {
					rookFrom := Position{
						X: move.From.X + 3,
						Y: move.From.Y,
					}
					rookTo := Position{
						X: move.From.X + 1, 
						Y: move.From.Y,
					}
					rookMove := CreateMove(rookFrom, rookTo)

					if ( !this.applyMoveBool(rookMove) ) { return false }
				}
			}
			piece.moveTo(move.To)
			this.Board.RemoveFrom(&move.From)
			this.Board.PlaceOn(piece, &move.To)
			this.LastMove = move
		}
		return true
	}
	return false
}

//func (s *State) ApplyMove(move *Move) {
	//piece := s.Board.GetPiece(&move.From)
	//if piece != nil {
		//piece.moveTo(move.To)
		//s.Board.RemoveFrom(&move.From)	
		//s.Board.PlaceOn(piece, &move.To)
		//s.LastMove = move
		//s.clearCache()
		//s.switchTurn()
	//} else {
		//println("piece was nil")
	//}
//}

func (s *State) ApplyMove(move *Move) (applied bool) {
	applied = s.applyMoveBool(move)

	if applied {
		s.switchTurn()
		s.clearCache()
	}
	return
}

func (this *State) isMoveLegal(move *Move) bool {
	if (len(this.legalMovesCache) != 0) {
		return slices.ContainsFunc(this.legalMovesCache, func(m *Move) bool {
			return move.Equal(m)
		}) 
	}

	next := this.Copy()
	next.ApplyMove(move)
	kingPos := next.Board.FindPiece("king", this.Turn)
	if (len(kingPos) == 0) {
		return false
	} 

	// condition 1
	controlledSquares := next.Board.squaresControlledBy(next.Turn);
	return !slices.ContainsFunc(controlledSquares, func(p *Position) bool {
		return p.Equal(*kingPos[0])	
	})
}

func (this *State) GetLegalMoves() []*Move {
	var cachedMoves []*Move
	if len(this.legalMovesCache) == 0 {
		cachedMoves = this.legalMovesCache
	} else {
		cachedMoves = this.legalMovesOrderedCache
	}
	if (len(cachedMoves) != 0) { return cachedMoves }


	possibleMoves := this.GetPossibleMoves()
	this.legalMovesCache = slices.DeleteFunc(possibleMoves, func(m *Move) bool {
		return m == nil || !this.isMoveLegal(m)
	})
	return this.legalMovesCache
}

func (this *State) IsGameOver() bool {
	return this.IsCheckmate() || this.IsStalemate()
}

func (this *State) IsCheckmate() bool {
	if legalMoves := this.GetLegalMoves(); len(legalMoves) == 0 {
		kingPos := this.Board.FindPiece("king", this.Turn)
		if len(kingPos) == 0 { return false }
		copy := this.Copy()
		copy.switchTurn()
		if slices.ContainsFunc(copy.GetPossibleMoves(), func(m *Move) bool {
			return m.To.Equal(*kingPos[0])
		}) {
			//fmt.Println("CHECKMATE")
			return true
		}
	}
	return false
}

func (this *State) IsStalemate() bool {
	if len(this.GetLegalMoves()) == 0 && !this.IsCheckmate() {
		//fmt.Println("STALEMATE")
		return true
	}
	return false
}


func (this *State) Copy() *State {
	newBoard := this.Board.copy()
	return &State{
		Board: newBoard,
		Turn: this.Turn,
		LastMove: this.LastMove,
		possibleMovesCache: nil,
		legalMovesCache: nil,
		legalMovesOrderedCache: nil,
	}
}

func (this *State) Equal(other *State) bool {
	return this.Board.Equal(other.Board) && this.Turn == other.Turn
}

func (this *State) String() string {
	return this.Board.String()
}

