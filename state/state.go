package state

import (
	"slices"
)

type State struct {
	Board *Board
	Turn string
	PreviousMoves []*Move
	possibleMovesCache []*Move
	legalMovesCache []*Move
	legalMovesOrderedCache []*Move
}

func CreateState(setup string) *State {
	if !(slices.Contains(setups, setup)) { setup = "clear" }
	return &State{
		Board: createBoard(setup),
		Turn: "white",
		PreviousMoves: []*Move{},
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


func (s *State) applyMoveBool(move *Move) bool {
	piece := s.Board.GetPiece(&move.From) 
	if piece != nil {
		if move.Promotion == "queen" {
			s.Board.RemoveFrom(&move.From)
			newPiece := s.Board.placeNew(piece.Color, move.Promotion, move.To)
			newPiece.HasMoved = true
			s.PreviousMoves = append(s.PreviousMoves, move) 
			return true
		}	else {
			if piece.Type == "king" {
				s.Board.clearKingsPosCash(piece.Color)
				if !piece.HasMoved {
					var rookMove *Move

					switch move.To.X - move.From.X {

					case 2: // short caslte
						rookMove = CreateMove(
							Position{
								X: move.From.X + 3,
								Y: move.From.Y,
							},
							Position{
								X: move.From.X + 1, 
								Y: move.From.Y,
							},
						)

					case -2: // long castle
						rookMove = CreateMove(
							Position{
								X: move.From.X - 4,
								Y: move.From.Y,
							},
							Position{
								X: move.From.X - 1,
								Y: move.From.Y,
							},
						)
					default:
						rookMove = nil
					}

					if rookMove == nil || !s.applyMoveBool(rookMove) { return false }
				}
			}
			piece.moveTo(move.To)
			s.Board.RemoveFrom(&move.From)
			s.Board.PlaceOn(piece, &move.To)
			s.PreviousMoves = append(s.PreviousMoves, move)
		}
		return true
	}
	return false
}

func (s *State) ApplyMove(move *Move) (applied bool) {
	applied = s.applyMoveBool(move)

	if applied {
		s.switchTurn()
		s.clearCache()
	}
	return
}


func (s *State) isMoveLegal(move *Move) bool {
	if (len(s.legalMovesCache) != 0) {
		return slices.ContainsFunc(s.legalMovesCache, func(m *Move) bool {
			return move.Equal(m)
		}) 
	}

	next := s.Copy()
	next.ApplyMove(move)
	kingPos := next.Board.FindPiece("king", s.Turn)
	if (len(kingPos) == 0) {
		return false
	} 

	// condition 1
	controlledSquares := next.Board.squaresControlledBy(next.Turn);
	return !slices.ContainsFunc(controlledSquares, func(p *Position) bool {
		return p.Equal(*kingPos[0])	
	})
}

func (s *State) GetLegalMoves() []*Move {
	var cachedMoves []*Move
	if len(s.legalMovesCache) == 0 {
		cachedMoves = s.legalMovesCache
	} else {
		cachedMoves = s.legalMovesOrderedCache
	}
	if (len(cachedMoves) != 0) { return cachedMoves }


	possibleMoves := s.GetPossibleMoves()
	s.legalMovesCache = slices.DeleteFunc(possibleMoves, func(m *Move) bool {
		return m == nil || !s.isMoveLegal(m)
	})
	return s.legalMovesCache
}


func (s *State) IsGameOver() bool {
	return s.IsCheckmate() || s.IsStalemate()
}

func (s *State) IsCheckmate() bool {
	if legalMoves := s.GetLegalMoves(); len(legalMoves) == 0 {
		kingPos := s.Board.FindPiece("king", s.Turn)
		if len(kingPos) == 0 || kingPos == nil { return false }
		copy := s.Copy()
		copy.switchTurn()
		if slices.ContainsFunc(copy.GetPossibleMoves(), func(m *Move) bool {
			return m.To.Equal(*kingPos[0])
		}) {
			return true
		}
	}
	return false
}

func (s *State) IsStalemate() bool {
	if len(s.GetLegalMoves()) == 0 && !s.IsCheckmate() {
		//fmt.Println("STALEMATE")
		return true
	}
	return false
}

func (s *State) Copy() *State {
	newBoard := s.Board.copy()
	return &State{
		Board: newBoard,
		Turn: s.Turn,
		PreviousMoves: s.PreviousMoves,
		possibleMovesCache: []*Move{},
		legalMovesCache: []*Move{},
		legalMovesOrderedCache: []*Move{},
	}
}

func (s *State) Equal(other *State) bool {
	return s.Board.Equal(other.Board) && s.Turn == other.Turn
}

func (s *State) String() string {
	return s.Board.String()
}

func (s *State) isCopy() bool {
	return s.Board.isCopy
}
