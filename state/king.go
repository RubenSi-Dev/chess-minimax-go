package state

var kingDirections = queenDirections


func (p Piece) getPossibleMovesKing(board *Board) (result []*Move) {
	result = p.stepMoves(board, &kingDirections)

	if p.isCastlingPossible(true, board) {
		result = append(result, CreateMove(
			p.Pos,
			Position{X: p.Pos.X + 2, Y: p.Pos.Y},
		))
	} 
	if p.isCastlingPossible(false, board) {
		result = append(result, CreateMove(
			p.Pos,
			Position{X: p.Pos.X - 2, Y: p.Pos.Y},
		))
	}
	return
}



func (this Piece) isCastlingPossible(short bool, board *Board) bool {
	if this.HasMoved {return false}
	posx := this.Pos.X
	posy := this.Pos.Y

	if short { // short castling
		rook := board.GetPiece(&Position{ X: posx + 3, Y: posy })

		if (
			board.GetPiece(&Position{ X: posx + 1, Y: posy }) != nil ||
			board.GetPiece(&Position{ X: posx + 2, Y: posy }) != nil ||
			rook == nil ||
			rook.Type != "rook" ||
			rook.HasMoved) {
			return false
		}
	} else { // long castling
		rook := board.GetPiece(&Position{ X: posx - 4, Y: posy })

		if (
			board.GetPiece(&Position{ X: posx - 1, Y: posy }) != nil ||
			board.GetPiece(&Position{ X: posx - 2, Y: posy }) != nil ||
			board.GetPiece(&Position{ X: posx - 3, Y: posy }) != nil ||
			rook == nil ||
			rook.Type != "rook" ||
			rook.HasMoved) {
			return false
		}
	}
	
	return true
}
