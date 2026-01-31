package state

import "fmt"

func CreateMovePawn(white bool, from *Position, to *Position) (*Move, error) {
	if from == nil || to == nil {
		return nil, fmt.Errorf("invalid position")
	}
	if white && to.Y == 7 {
		return CreateMovePromotion(*from, *to, "queen"), nil
	} else if !white && to.Y == 0 {
		return CreateMovePromotion(*from, *to, "queen"), nil
	}
	return CreateMove(*from, *to), nil
}

func (p *Piece) getPossibleMovesPawn(board *Board) ([]*Move, error) {
	possibleMoves := []*Move{}
	white := (p.Color == "white")

	var move0 Position
	if white {
		move0 = Position{X: p.Pos.X, Y: p.Pos.Y + 1}
	} else {
		move0 = Position{X: p.Pos.X, Y: p.Pos.Y - 1}
	}

	target, err := board.GetPiece(&move0)
	if err != nil {
		return nil, err
	}
	inBounds := board.isInBounds(&move0)
	if target == nil && inBounds {
		pawnMove, err := CreateMovePawn(white, &p.Pos, &move0)
		if err != nil {
			return nil, err
		}
		possibleMoves = append(possibleMoves, pawnMove)
		if (white && p.Pos.Y == 1) || (!white && p.Pos.Y == 6) {
			var move1 Position
			if white {
				move1 = Position{X: p.Pos.X, Y: p.Pos.Y + 2}
			} else {
				move1 = Position{X: p.Pos.X, Y: p.Pos.Y - 2}
			}
			if res, err := board.GetPiece(&move1); res == nil {
				if err != nil {
					return nil, err
				}
				pawnMove, err := CreateMovePawn(white, &p.Pos, &move1)
				if err != nil {
					return nil, err
				}
				possibleMoves = append(possibleMoves, pawnMove)
			}
		}
	}

	var movesDiag []Position

	if white {
		movesDiag = []Position{
			{X: p.Pos.X + 1, Y: p.Pos.Y + 1},
			{X: p.Pos.X - 1, Y: p.Pos.Y + 1}}
	} else {
		movesDiag = []Position{
			{X: p.Pos.X + 1, Y: p.Pos.Y - 1},
			{X: p.Pos.X - 1, Y: p.Pos.Y - 1}}
	}

	for _, move := range movesDiag {
		if board.isInBounds(&move) {
			target, err := board.GetPiece(&move)
			if err != nil {
				return nil, err
			}
			if target != nil && target.Color != p.Color {
				pawnmove, err := CreateMovePawn(white, &p.Pos, &move)
				if err != nil {
					return nil, err
				}
				possibleMoves = append(possibleMoves, pawnmove)
			}
		}
	}
	return possibleMoves, nil
}
