package state

var queenDirections = append(rookDirections, bishopDirections...)

func (p Piece) getPossibleMovesQueen(board *Board) ([]*Move, error) {
	raymoves, err := p.rayMoves(board, &queenDirections)
	if err != nil {
		return nil, err
	}
	return raymoves, nil
}
