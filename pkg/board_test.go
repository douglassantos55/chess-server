package pkg

import (
	"reflect"
	"testing"
)

func TestParseSquare(t *testing.T) {
	if s, _ := parseSquare("a1"); s.col != 'a' {
		t.Errorf("Expected 'a', got '%v'", s)
	}
	if s, _ := parseSquare("a1"); s.row != 0 {
		t.Errorf("Expected '0', got '%v'", s)
	}

	if s, _ := parseSquare("d5"); s.col != 'd' {
		t.Errorf("Expected 'd', got '%v'", s)
	}
	if s, _ := parseSquare("d5"); s.row != 4 {
		t.Errorf("Expected '4', got '%v'", s)
	}

	if s, _ := parseSquare("h8"); s.col != 'h' {
		t.Errorf("Expected 'h', got '%v'", s)
	}
	if s, _ := parseSquare("h8"); s.row != 7 {
		t.Errorf("Expected '7', got '%v'", s)
	}
}

func TestCreatesBoard(t *testing.T) {
	board := NewBoard()

	if board.Square("a1") != Rook(White) {
		t.Errorf("Expected Rook on a1, got %v", board.Square("a1"))
	}
	if board.Square("b1") != Knight(White) {
		t.Errorf("Expected Knight on b1, got %v", board.Square("b1"))
	}
	if board.Square("c1") != Bishop(White) {
		t.Errorf("Expected Bishop on c1, got %v", board.Square("c1"))
	}
	if !reflect.DeepEqual(board.Square("d1"), Queen(White)) {
		t.Errorf("Expected Queen on d1, got %v", board.Square("d1"))
	}
	if !reflect.DeepEqual(board.Square("e1"), King(White)) {
		t.Errorf("Expected King on e1, got %v", board.Square("e1"))
	}
	if board.Square("f1") != Bishop(White) {
		t.Errorf("Expected Bishop on f1, got %v", board.Square("f1"))
	}
	if board.Square("g1") != Knight(White) {
		t.Errorf("Expected Knight on g1, got %v", board.Square("g1"))
	}
	if board.Square("h1") != Rook(White) {
		t.Errorf("Expected Rook on h1, got %v", board.Square("h1"))
	}

	if board.Square("a8") != Rook(Black) {
		t.Errorf("Expected Rook on a8, got %v", board.Square("a8"))
	}
	if board.Square("b8") != Knight(Black) {
		t.Errorf("Expected Knight on b8, got %v", board.Square("b8"))
	}
	if board.Square("c8") != Bishop(Black) {
		t.Errorf("Expected Bishop on c8, got %v", board.Square("c8"))
	}
	if !reflect.DeepEqual(board.Square("d8"), Queen(Black)) {
		t.Errorf("Expected Queen on d8, got %v", board.Square("d8"))
	}
	if !reflect.DeepEqual(board.Square("e8"), King(Black)) {
		t.Errorf("Expected King on e8, got %v", board.Square("e8"))
	}
	if board.Square("f8") != Bishop(Black) {
		t.Errorf("Expected Bishop on f8, got %v", board.Square("f8"))
	}
	if board.Square("g8") != Knight(Black) {
		t.Errorf("Expected Knight on g8, got %v", board.Square("g8"))
	}
	if board.Square("h8") != Rook(Black) {
		t.Errorf("Expected Rook on h8, got %v", board.Square("h8"))
	}
}

func TestMovePiece(t *testing.T) {
	board := NewBoard()
	board.Move("b1", "c3")

	if board.Square("b1") == Knight(White) {
		t.Errorf("Expected empty square on b1, got %v", board.Square("b1"))
	}
	if board.Square("c3") != Knight(White) {
		t.Errorf("Expected Knight on c3, got %v", board.Square("c3"))
	}

	board.Move("c3", "b1")

	if board.Square("b1") != Knight(White) {
		t.Errorf("Expected Knight on b1, got %v", board.Square("b1"))
	}
	if board.Square("c3") == Knight(White) {
		t.Errorf("Expected empty square on c3, got %v", board.Square("c3"))
	}

	board.Move("b1", "c3")
	board.Move("c3", "b5")

	if board.Square("c3") == Knight(White) {
		t.Errorf("Expected empty square on c3, got %v", board.Square("c3"))
	}
	if board.Square("b5") != Knight(White) {
		t.Errorf("Expected Knight on b5, got %v", board.Square("b5"))
	}

	board.Move("a2", "c3")

	if board.Square("c3") == Pawn(White) {
		t.Error("Should not have a pawn on c3, invalid movement")
	}

	board.Move("a1", "a5")

	if board.Square("a5") == Rook(White) {
		t.Error("Should not move from a1 to a5, there's a pawn in the way")
	}

	board.Move("a2", "a4")
	board.Move("a1", "a4")

	if board.Square("a4") == Rook(White) {
		t.Error("Should not move from a1 to a4, there's a pawn in the way")
	}

	board.Move("a1", "a3")

	if board.Square("a3") != Rook(White) {
		t.Error("Should move Rook from a1 to a3")
	}
}

func TestMoveOverOpponentsPieces(t *testing.T) {
	board := NewBoard()

	board.Move("d2", "d4")

	if board.Square("d2") == Pawn(White) {
		t.Errorf("Expected empty square on d2, got %v", board.Square("d2"))
	}
	if board.Square("d4") != Pawn(White) {
		t.Errorf("Expected pawn on d4, got %v", board.Square("d4"))
	}

	board.Move("e7", "e5")

	if board.Square("e7") == Pawn(Black) {
		t.Errorf("Expected empty square on e7, got %v", board.Square("e7"))
	}
	if board.Square("e5") != Pawn(Black) {
		t.Errorf("Expected pawn on e5, got %v", board.Square("e5"))
	}

	board.Move("f8", "c5")

	if board.Square("c5") != Bishop(Black) {
		t.Errorf("Expected bishop on c5, got %v", board.Square("c5"))
	}

	board.Move("c5", "e3")

	if board.Square("e3") == Bishop(Black) {
		t.Errorf("Should not jump over d4")
	}
}

func TestForwardCaptureDiagonally(t *testing.T) {
	board := NewBoard()

	board.Move("e2", "e4")
	board.Move("e7", "e5")
	board.Move("d2", "d4")
	board.Move("e5", "d4")

	if board.Square("d4") != Pawn(Black) {
		t.Error("Should capture d4 from e5")
	}

	board.Move("e4", "f5")
	if board.Square("f5") != Empty() {
		t.Error("Should not capture empty square")
	}

	board.Move("e4", "d5")
	if board.Square("d5") != Empty() {
		t.Error("Should not capture empty square")
	}

	board.Move("e4", "d3")
	if board.Square("d3") != Empty() {
		t.Error("Should not capture backwards")
	}

	board.Move("e4", "f3")
	if board.Square("f3") != Empty() {
		t.Error("Should not capture backwards")
	}

	board.Move("e4", "e5")
	board.Move("e5", "d4")
	if board.Square("d4") != Pawn(Black) {
		t.Error("Should not capture backwards")
	}

	board.Move("d4", "e5")
	if board.Square("e5") != Pawn(White) {
		t.Error("Should not capture backwards")
	}

	board.Move("e5", "e6")
	board.Move("e6", "d7")
	if board.Square("d7") != Pawn(White) {
		t.Error("Should capture d7 from e6")
	}

	board.Move("d7", "e8")
	if board.Square("e8") != Pawn(White) {
		t.Error("Should capture e8 from d7")
	}

	board.Move("d4", "d3")
	board.Move("d3", "e2")
	if board.Square("e2") != Empty() {
		t.Error("Should not capture empty square")
	}

	board.Move("d3", "c2")
	if board.Square("c2") != Pawn(Black) {
		t.Error("Should capture c2 from d3")
	}

	board.Move("c2", "d1")
	if board.Square("d1") != Pawn(Black) {
		t.Error("Should capture d1 from c2")
	}
}
