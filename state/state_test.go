package state

import (
	"testing"
)

func TestCreateState(t *testing.T) {
	s, err := CreateState("default")
	if err != nil {
		t.Fatalf("Failed to create state: %v", err)
	}
	if s.Turn != "white" {
		t.Errorf("Expected turn to be white, got %s", s.Turn)
	}

	// Check board setup
	// White Rook at A1 (0,0)
	p, err := s.Board.GetPiece(&Position{X: 0, Y: 0})
	if err != nil {
		t.Fatalf("Error getting piece: %v", err)
	}
	if p == nil {
		t.Fatalf("Expected piece at 0,0")
	}
	if p.Type != "rook" || p.Color != "white" {
		t.Errorf("Expected white rook at 0,0, got %s %s", p.Color, p.Type)
	}

	// Black King at E8 (4,7)
	p, err = s.Board.GetPiece(&Position{X: 4, Y: 7})
	if err != nil {
		t.Fatalf("Error getting piece: %v", err)
	}
	if p == nil {
		t.Fatalf("Expected piece at 4,7")
	}
	if p.Type != "king" || p.Color != "black" {
		t.Errorf("Expected black king at 4,7, got %s %s", p.Color, p.Type)
	}
}

func TestPawnMoves(t *testing.T) {
	s, err := CreateState("default")
	if err != nil {
		t.Fatal(err)
	}
	// White pawn at E2 (4, 1)
	pawnPos := Position{X: 4, Y: 1}
	pawn, _ := s.Board.GetPiece(&pawnPos)
	moves, err := pawn.GetPossibleMoves(s.Board)
	if err != nil {
		t.Fatal(err)
	}
	// Should have 2 moves: E3 (4,2), E4 (4,3)
	if len(moves) != 2 {
		t.Errorf("Expected 2 moves for pawn at start, got %d", len(moves))
	}

	foundE3 := false
	foundE4 := false
	for _, m := range moves {
		if m.To.X == 4 && m.To.Y == 2 {
			foundE3 = true
		}
		if m.To.X == 4 && m.To.Y == 3 {
			foundE4 = true
		}
	}
	if !foundE3 || !foundE4 {
		t.Errorf("Expected moves to E3 and E4, got %v", moves)
	}
}

func TestApplyMove(t *testing.T) {
	s, err := CreateState("default")
	if err != nil {
		t.Fatal(err)
	}
	// Move E2 to E4
	move := CreateMove(Position{X: 4, Y: 1}, Position{X: 4, Y: 3})
	applied, err := s.ApplyMove(move)
	if err != nil {
		t.Fatal(err)
	}
	if !applied {
		t.Fatal("Move should be applied")
	}
	if s.Turn != "black" {
		t.Errorf("Turn should switch to black")
	}
	p, _ := s.Board.GetPiece(&Position{X: 4, Y: 3})
	if p == nil || p.Type != "pawn" {
		t.Errorf("Pawn should be at E4")
	}
	oldP, _ := s.Board.GetPiece(&Position{X: 4, Y: 1})
	if oldP != nil {
		t.Errorf("E2 should be empty")
	}
}

func TestFoolsMate(t *testing.T) {
	s, err := CreateState("default")
	if err != nil {
		t.Fatal(err)
	}

	// Fool's Mate sequence
	moves := []string{
		"F2-F3", // White
		"E7-E5", // Black
		"G2-G4", // White
		"D8-H4", // Black mate
	}

	for _, alg := range moves {
		m := FromAlgebraicToMove(alg)
		_, err := s.ApplyMove(m)
		if err != nil {
			t.Fatalf("Failed to apply move %s: %v", alg, err)
		}
	}

	mate, err := s.IsCheckmate()
	if err != nil {
		t.Fatal(err)
	}
	if !mate {
		t.Error("Expected checkmate")
	}

	gameOver, err := s.IsGameOver()
	if err != nil {
		t.Fatal(err)
	}
	if !gameOver {
		t.Error("Expected game over")
	}
}

func TestCastling(t *testing.T) {
	s, err := CreateState("castling")
	if err != nil {
		t.Fatal(err)
	}

	// Check possible moves for White King at E1 (4,0)
	kingPos := Position{X: 4, Y: 0}
	king, _ := s.Board.GetPiece(&kingPos)
	moves, err := king.GetPossibleMoves(s.Board)
	if err != nil {
		t.Fatal(err)
	}

	// Should include castling moves
	hasShort := false
	hasLong := false
	for _, m := range moves {
		if m.To.X == 6 && m.To.Y == 0 {
			hasShort = true
		}
		if m.To.X == 2 && m.To.Y == 0 {
			hasLong = true
		}
	}

	if !hasShort {
		t.Error("Short castling should be possible")
	}
	if !hasLong {
		t.Error("Long castling should be possible")
	}

	// Perform short castle
	castleMove := CreateMove(kingPos, Position{X: 6, Y: 0})
	_, err = s.ApplyMove(castleMove)
	if err != nil {
		t.Fatal(err)
	}

	// Check King position
	k, _ := s.Board.GetPiece(&Position{X: 6, Y: 0})
	if k == nil || k.Type != "king" {
		t.Error("King should be at G1")
	}

	// Check Rook position (H1 -> F1)
	r, _ := s.Board.GetPiece(&Position{X: 5, Y: 0})
	if r == nil || r.Type != "rook" {
		t.Error("Rook should be at F1 after short castle")
	}
}

func TestPromotion(t *testing.T) {
	s, err := CreateState("promotion")
	if err != nil {
		t.Fatal(err)
	}

	// White pawn at D7 (3, 6) ready to promote
	// Move to D8 (3, 7)
	from := Position{X: 3, Y: 6}
	to := Position{X: 3, Y: 7}

	// Check if move is possible
	pawn, _ := s.Board.GetPiece(&from)
	moves, _ := pawn.GetPossibleMoves(s.Board)

	foundPromotion := false
	for _, m := range moves {
		if m.To.Equal(to) && m.Promotion == "queen" {
			foundPromotion = true
		}
	}
	if !foundPromotion {
		t.Error("Expected promotion move to queen")
	}

	// Apply promotion
	move := CreateMovePromotion(from, to, "queen")
	_, err = s.ApplyMove(move)
	if err != nil {
		t.Fatal(err)
	}

	newPiece, _ := s.Board.GetPiece(&to)
	if newPiece == nil || newPiece.Type != "queen" {
		t.Error("Pawn should have promoted to queen")
	}
}
