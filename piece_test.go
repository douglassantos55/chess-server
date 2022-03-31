package main

import "testing"

func TestStraightMovement(t *testing.T) {
	straight := Straight{}

	if !straight.IsValid("a1", "h1") {
		t.Error("Should move from a1 to h1")
	}
	if straight.IsValid("a1", "b2") {
		t.Error("Should not move from a1 to b2")
	}

	straight.squares = 2
	if straight.IsValid("a1", "a7") {
		t.Error("Should not move from a1 to a7")
	}
	if straight.IsValid("a1", "f1") {
		t.Error("Should not move from a1 to f1")
	}
}

func TestDiagonalMovement(t *testing.T) {
	diagonal := Diagonal{}

	if !diagonal.IsValid("a1", "b2") {
		t.Error("Should move from a1 to b2")
	}
	if !diagonal.IsValid("a1", "h8") {
		t.Error("Should move from a1 to h8")
	}
	if !diagonal.IsValid("h8", "a1") {
		t.Error("Should move from h8 to a1")
	}
	if !diagonal.IsValid("a8", "h1") {
		t.Error("Should move from a8 to h1")
	}
	if diagonal.IsValid("a1", "b1") {
		t.Error("Should not move from a1 to b1")
	}
	if diagonal.IsValid("a1", "c1") {
		t.Error("Should not move from a1 to c1")
	}
	if diagonal.IsValid("a1", "a8") {
		t.Error("Should not move from a1 to a8")
	}
	if diagonal.IsValid("b1", "c3") {
		t.Error("Should not move from b1 to c3")
	}
	if diagonal.IsValid("a1", "b3") {
		t.Error("Should not move from a1 to b3")
	}

	diagonal.squares = 1
	if !diagonal.IsValid("a1", "b2") {
		t.Error("Should move from a1 to b2")
	}
	if !diagonal.IsValid("h8", "g7") {
		t.Error("Should move from h8 to g7")
	}
	if diagonal.IsValid("a1", "c3") {
		t.Error("Should not move from a1 to c3")
	}
	if diagonal.IsValid("c3", "a1") {
		t.Error("Should not move from c3 to a1")
	}
}

func TestLMovement(t *testing.T) {
	l := LMovement{}

	if !l.IsValid("b1", "c3") {
		t.Error("Should move from b1 to c3")
	}
	if !l.IsValid("b1", "a3") {
		t.Error("Should move from b1 to a3")
	}
	if !l.IsValid("c3", "b1") {
		t.Error("Should move from c3 to b1")
	}
	if l.IsValid("b1", "c2") {
		t.Error("Should not move from b1 to c2")
	}
}

func TestForwardMovement(t *testing.T) {
	forward := Forward{1}

	if !forward.IsValid("a1", "a2") {
		t.Error("Should move from a1 to a2")
	}
	if !forward.IsValid("a2", "a3") {
		t.Error("Should move from a2 to a3")
	}
	if !forward.IsValid("h1", "h2") {
		t.Error("Should move from h1 to h2")
	}
	if forward.IsValid("a1", "b1") {
		t.Error("Should not move from a1 to b1")
	}
	if forward.IsValid("a1", "a3") {
		t.Error("Should not move from a1 to a3")
	}
	if forward.IsValid("a2", "a1") {
		t.Error("Should not move backwards from a2 to a1")
	}
	if forward.IsValid("a2", "b3") {
		t.Error("Should not move from a2 to b3")
	}
}

func TestCombinedMovement(t *testing.T) {
	queen := Combined{
		[]Movement{Straight{}, Diagonal{}},
	}

	if !queen.IsValid("a1", "h8") {
		t.Error("Should move from a1 to h8")
	}
	if !queen.IsValid("a1", "a8") {
		t.Error("Should move from a1 to a8")
	}
	if !queen.IsValid("a1", "h1") {
		t.Error("Should move from a1 to h1")
	}
	if queen.IsValid("b1", "c3") {
		t.Error("Should not move from b1 to c3")
	}

	king := Combined{
		[]Movement{Straight{1}, Diagonal{1}},
	}

	if !king.IsValid("a1", "b2") {
		t.Error("Should move from a1 to b2")
	}
	if !king.IsValid("b2", "a1") {
		t.Error("Should move from b2 to a1")
	}
	if !king.IsValid("b2", "b3") {
		t.Error("Should move from b2 to b3")
	}
	if !king.IsValid("b2", "c1") {
		t.Error("Should move from b2 to c1")
	}
	if !king.IsValid("b2", "b1") {
		t.Error("Should move from b2 to b1")
	}
	if king.IsValid("b2", "b7") {
		t.Error("Should not move from b2 to b7")
	}
	if king.IsValid("b2", "d4") {
		t.Error("Should not move from b2 to d4")
	}
}

func TestStraightIsAllowed(t *testing.T) {
	straight := Straight{}

	if straight.IsAllowed("a1", "a5", White, NewBoard()) {
		t.Error("Should not jump over a2")
	}
	if straight.IsAllowed("a1", "h1", White, NewBoard()) {
		t.Error("Should not jump over b1")
	}
	if straight.IsAllowed("a5", "a1", White, NewBoard()) {
		t.Error("Should not jump over a2")
	}
	if straight.IsAllowed("h1", "a1", White, NewBoard()) {
		t.Error("Should not jump over g1")
	}
	if !straight.IsAllowed("c3", "c6", White, NewBoard()) {
		t.Error("Should be able to move from c3 to c6")
	}
	if !straight.IsAllowed("c3", "h3", White, NewBoard()) {
		t.Error("Should be able to move from c3 to h3")
	}
	if !straight.IsAllowed("c6", "c3", White, NewBoard()) {
		t.Error("Should be able to move from c6 to c3")
	}
	if !straight.IsAllowed("h3", "c3", White, NewBoard()) {
		t.Error("Should be able to move from h3 to c3")
	}
	if straight.IsAllowed("a3", "a8", White, NewBoard()) {
		t.Error("Should not be able to capture Rook at a8")
	}
	if !straight.IsAllowed("a3", "a7", White, NewBoard()) {
		t.Error("Should be able to capture pawn at a7")
	}
}

func TestDiagonalIsAllowed(t *testing.T) {
	diagonal := Diagonal{}

	if diagonal.IsAllowed("c1", "e3", White, NewBoard()) {
		t.Error("Should not jump over d2")
	}
	if diagonal.IsAllowed("e3", "c1", White, NewBoard()) {
		t.Error("Should not jump over d2")
	}
	if !diagonal.IsAllowed("e3", "h6", White, NewBoard()) {
		t.Error("Should be able to move from e3 to h6")
	}
	if !diagonal.IsAllowed("h6", "e3", White, NewBoard()) {
		t.Error("Should be able to move from h6 to e3")
	}
	if !diagonal.IsAllowed("c3", "g7", White, NewBoard()) {
		t.Error("Should be able to move from c3 to g7")
	}
	if diagonal.IsAllowed("c3", "h8", White, NewBoard()) {
		t.Error("Should not be able to capture from c3 at h8")
	}
	if diagonal.IsAllowed("f1", "c4", White, NewBoard()) {
		t.Error("Should not be able to move from f1 to c4")
	}
	if !diagonal.IsAllowed("c3", "g7", White, NewBoard()) {
		t.Error("Should be able to move from c3 to g7")
	}
	if diagonal.IsAllowed("c4", "f1", White, NewBoard()) {
		t.Error("Should not be able to move from c4 to f1")
	}
	if !diagonal.IsAllowed("c4", "e6", White, NewBoard()) {
		t.Error("Should not be able to move from c4 to e6")
	}
	if !diagonal.IsAllowed("e6", "c4", White, NewBoard()) {
		t.Error("Should not be able to move from e6 to c4")
	}
}

func TestForwardIsAllowed(t *testing.T) {
	forward := Forward{}

	if forward.IsAllowed("a1", "a2", White, NewBoard()) {
		t.Error("Should not be able to move from a1 to a2")
	}
	if forward.IsAllowed("a6", "a7", White, NewBoard()) {
		t.Error("Should not be able to move from a6 to a7")
	}
	if !forward.IsAllowed("a3", "a4", White, NewBoard()) {
		t.Error("Should be able to move from a3 to a4")
	}
}

func TestLMovementIsAllowed(t *testing.T) {
	l := LMovement{}

	if !l.IsAllowed("b1", "c3", White, NewBoard()) {
		t.Error("Should be able to move from b1 to c3")
	}
	if !l.IsAllowed("b1", "a3", White, NewBoard()) {
		t.Error("Should be able to move from b1 to a3")
	}
	if l.IsAllowed("c4", "d2", White, NewBoard()) {
		t.Error("Should not be able to move from c4 to d2")
	}
	if !l.IsAllowed("c5", "b7", White, NewBoard()) {
		t.Error("Should be able to capture from c5 at b7")
	}
	if !l.IsAllowed("c6", "d8", White, NewBoard()) {
		t.Error("Should be able to capture from c6 at d8")
	}
}

func TestCombinedIsAllowed(t *testing.T) {
	combined := Combined{
		[]Movement{Straight{}, Diagonal{}},
	}

	if combined.IsAllowed("a1", "a6", White, NewBoard()) {
		t.Error("Should not be able to move from a1 to a6")
	}
	if combined.IsAllowed("a6", "a2", White, NewBoard()) {
		t.Error("Should not be able to move from a6 to a1")
	}
	if combined.IsAllowed("a3", "a8", White, NewBoard()) {
		t.Error("Should not be able to move from a3 to a8")
	}
	if !combined.IsAllowed("a3", "f6", White, NewBoard()) {
		t.Error("Should be able to move from a3 to f6")
	}
	if !combined.IsAllowed("f6", "a3", White, NewBoard()) {
		t.Error("Should be able to move from f6 to a3")
	}
	if !combined.IsAllowed("c3", "g7", White, NewBoard()) {
		t.Error("Should be able to move from c3 to g7")
	}
	if combined.IsAllowed("c3", "h8", White, NewBoard()) {
		t.Error("Should not be able to capture from c3 at h8")
	}
}
