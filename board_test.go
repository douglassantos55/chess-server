package main

import "testing"

func TestParseSquare(t *testing.T) {
	board := NewBoard()

	if _, c := board.parseSquare("a1"); c != "a" {
		t.Errorf("Expected 'a', got '%v'", c)
	}
	if r, _ := board.parseSquare("a1"); r != 0 {
		t.Errorf("Expected '0', got '%v'", r)
	}

	if _, c := board.parseSquare("d5"); c != "d" {
		t.Errorf("Expected 'd', got '%v'", c)
	}
	if r, _ := board.parseSquare("d5"); r != 4 {
		t.Errorf("Expected '4', got '%v'", r)
	}

	if _, c := board.parseSquare("h8"); c != "h" {
		t.Errorf("Expected 'h', got '%v'", c)
	}
	if r, _ := board.parseSquare("h8"); r != 7 {
		t.Errorf("Expected '7', got '%v'", r)
	}
}

func TestCreatesBoard(t *testing.T) {
	board := NewBoard()

	if board.Square("a1") != "R" {
		t.Errorf("Expected Rook on a1, got %v", board.Square("a1"))
	}
	if board.Square("b1") != "N" {
		t.Errorf("Expected Knight on b1, got %v", board.Square("b1"))
	}
	if board.Square("c1") != "B" {
		t.Errorf("Expected Bishop on c1, got %v", board.Square("c1"))
	}
	if board.Square("d1") != "Q" {
		t.Errorf("Expected Queen on d1, got %v", board.Square("d1"))
	}
	if board.Square("e1") != "K" {
		t.Errorf("Expected King on e1, got %v", board.Square("e1"))
	}
	if board.Square("f1") != "B" {
		t.Errorf("Expected Bishop on f1, got %v", board.Square("f1"))
	}
	if board.Square("g1") != "N" {
		t.Errorf("Expected Knight on g1, got %v", board.Square("g1"))
	}
	if board.Square("h1") != "R" {
		t.Errorf("Expected Rook on h1, got %v", board.Square("h1"))
	}

	if board.Square("a8") != "R" {
		t.Errorf("Expected Rook on a8, got %v", board.Square("a8"))
	}
	if board.Square("b8") != "N" {
		t.Errorf("Expected Knight on b8, got %v", board.Square("b8"))
	}
	if board.Square("c8") != "B" {
		t.Errorf("Expected Bishop on c8, got %v", board.Square("c8"))
	}
	if board.Square("d8") != "Q" {
		t.Errorf("Expected Queen on d8, got %v", board.Square("d8"))
	}
	if board.Square("e8") != "K" {
		t.Errorf("Expected King on e8, got %v", board.Square("e8"))
	}
	if board.Square("f8") != "B" {
		t.Errorf("Expected Bishop on f8, got %v", board.Square("f8"))
	}
	if board.Square("g8") != "N" {
		t.Errorf("Expected Knight on g8, got %v", board.Square("g8"))
	}
	if board.Square("h8") != "R" {
		t.Errorf("Expected Rook on h8, got %v", board.Square("h8"))
	}
}

func TestMovePiece(t *testing.T) {
	board := NewBoard()
	board.Move("b1", "c3")

	if board.Square("b1") == "N" {
		t.Errorf("Expected empty square on b1, got %v", board.Square("b1"))
	}
	if board.Square("c3") != "N" {
		t.Errorf("Expected Knight on c3, got %v", board.Square("c3"))
	}

	board.Move("c3", "b5")

	if board.Square("c3") == "N" {
		t.Errorf("Expected empty square on c3, got %v", board.Square("c3"))
	}
	if board.Square("b5") != "N" {
		t.Errorf("Expected Knight on b5, got %v", board.Square("b5"))
	}
}