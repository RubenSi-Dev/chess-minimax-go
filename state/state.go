package state

import (
	"slices"
)

// State - represents the current state of a chess game
type State struct {
	Board                  *Board // Board representation
	Turn                   string // "white" or "black"
	PreviousMoves          []*Move
	possibleMovesCache     []*Move
	legalMovesCache        []*Move
	legalMovesOrderedCache []*Move
}

// CreateState - creates a new game state with the given setup
func CreateState(setup string) (*State, error) {
	if !(slices.Contains(setups, setup)) {
		setup = "clear"
	}
	newBoard, err := createBoard(setup)
	if err != nil {
		return nil, err
	}
	return &State{
		Board:                  newBoard,
		Turn:                   "white",
		PreviousMoves:          []*Move{},
		possibleMovesCache:     []*Move{},
		legalMovesCache:        []*Move{},
		legalMovesOrderedCache: []*Move{},
	}, nil
}

// clearCache - clears cached possible and legal moves
// usually called when the state changes (a move is played)
func (s *State) clearCache() {
	s.possibleMovesCache = []*Move{}
	s.legalMovesCache = []*Move{}
	s.legalMovesOrderedCache = []*Move{}
	s.Board.clearCache()
}

func (s *State) switchTurn() string {
	white := s.Turn == "white"
	if white {
		s.Turn = "black"
	} else {
		s.Turn = "white"
	}
	return s.Turn
}

// GetPossibleMoves - get all possible moves in the current state (inoring illegal moves by board context)
// ranges over board.GetPieces() and appends it to the result (saves it in cache too)
// first checks whether its already in cache
func (s *State) GetPossibleMoves() []*Move {
	if len(s.possibleMovesCache) != 0 {
		return s.possibleMovesCache
	}

	pieces := s.Board.GetPieces()
	for _, piece := range pieces {
		if piece.Color == s.Turn {
			moves := piece.GetPossibleMoves(s.Board)
			s.possibleMovesCache = append(s.possibleMovesCache, moves...)
		}
	}
	return s.possibleMovesCache
}

// applyMoveBool - helper function that returns true if the move is possible
// helpful for moves like castling where two pieces have to move
func (s *State) applyMoveBool(move *Move) (bool, error) {
	piece := s.Board.GetPiece(&move.From)
	if piece != nil {
		if move.Promotion == "queen" {
			s.Board.RemoveFrom(&move.From)
			newPiece, err := s.Board.placeNew(piece.Color, move.Promotion, move.To)
			if err != nil {
				return false, err
			}
			newPiece.HasMoved = true
			s.PreviousMoves = append(s.PreviousMoves, move)
			return true, nil
		} else {
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

					boolApply, err := s.applyMoveBool(rookMove)
					if rookMove == nil || err != nil || !boolApply {
						return false, err
					}
				}
			}
			piece.moveTo(move.To)
			s.Board.RemoveFrom(&move.From)
			s.Board.PlaceOn(piece, &move.To)
			s.PreviousMoves = append(s.PreviousMoves, move)
		}
		return true, nil
	}
	return false, nil
}

// ApplyMove - wrapper function for applyMoveBool, then switches the turn
func (s *State) ApplyMove(move *Move) (bool, error) {
	applied, err := s.applyMoveBool(move)
	if err != nil {
		return false, err
	}

	if applied {
		s.switchTurn()
		s.clearCache()
	}
	return applied, nil
}

// isMoveLegal - checks whether a *possiblemove* is also *legal*
// its legal if your own king is not in check after the move is completed
// required to make a deepcopy of this state, play the move, and see if there are problems
func (s *State) isMoveLegal(move *Move) (bool, error) {
	if len(s.legalMovesCache) != 0 {
		return slices.ContainsFunc(s.legalMovesCache, func(m *Move) bool {
			return move.Equal(m)
		}), nil
	}

	next, err := s.Copy()
	if err != nil {
		return false, err
	}

	_, err = next.ApplyMove(move)
	if err != nil {
		return false, err
	}
	kingPos := next.Board.FindPiece("king", s.Turn)
	if len(kingPos) == 0 {
		return false, nil
	}

	controlledSquares := next.Board.squaresControlledBy(next.Turn)
	return !slices.ContainsFunc(controlledSquares, func(p *Position) bool {
		return p.Equal(*kingPos[0])
	}), nil
}

func (s *State) GetLegalMoves() []*Move {
	var cachedMoves []*Move
	if len(s.legalMovesCache) == 0 {
		cachedMoves = s.legalMovesCache
	} else {
		cachedMoves = s.legalMovesOrderedCache
	}
	if len(cachedMoves) != 0 {
		return cachedMoves
	}

	possibleMoves := s.GetPossibleMoves()
	s.legalMovesCache = slices.DeleteFunc(possibleMoves, func(m *Move) bool {
		isLegal, _ := s.isMoveLegal(m)
		return m == nil || !isLegal
	})
	return s.legalMovesCache
}

func (s *State) IsGameOver() (bool, error) {
	isMate, err := s.IsCheckmate()
	if err != nil {
		return false, err
	}

	isStale, err := s.IsStalemate()
	if err != nil {
		return false, err
	}

	return (isMate || isStale), nil
}

// IsCheckmate
// check by switching turns and seeing if the player has any legal moves
// also checks for stalemate
func (s *State) IsCheckmate() (bool, error) {
	if legalMoves := s.GetLegalMoves(); len(legalMoves) == 0 {
		kingPos := s.Board.FindPiece("king", s.Turn)
		if len(kingPos) == 0 || kingPos == nil {
			return false, nil
		}
		copy, err := s.Copy()
		if err != nil {
			return false, err
		}
		copy.switchTurn()
		if slices.ContainsFunc(copy.GetPossibleMoves(), func(m *Move) bool {
			return m.To.Equal(*kingPos[0])
		}) {
			return true, nil
		}
	}
	return false, nil
}

func (s *State) IsStalemate() (bool, error) {
	isMate, err := s.IsCheckmate()
	if err != nil {
		return false, err
	}

	if len(s.GetLegalMoves()) == 0 && !isMate {
		//fmt.Println("STALEMATE")
		return true, nil
	}
	return false, nil
}

// Copy - makes a deep copy of the state
// deepcopy -> all pieces on the board most be copied and so on
func (s *State) Copy() (*State, error) {
	newBoard, err := s.Board.copy()
	if err != nil {
		return nil, err
	}
	return &State{
		Board:                  newBoard,
		Turn:                   s.Turn,
		PreviousMoves:          []*Move{},
		possibleMovesCache:     []*Move{},
		legalMovesCache:        []*Move{},
		legalMovesOrderedCache: []*Move{},
	}, nil
}

// some helper functions
func (s *State) Equal(other *State) bool {
	return s.Board.Equal(other.Board) && s.Turn == other.Turn
}

func (s *State) String() string {
	return s.Board.String()
}

func (s *State) isCopy() bool {
	return s.Board.isCopy
}
