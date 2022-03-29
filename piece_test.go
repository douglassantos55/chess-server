package main

import "testing"

func TestStraightMovement(t *testing.T) {
    straight := Straight{}
    if !straight.Move("a1", "h1") {
        t.Error("Should move from a1 to h1")
    }
    if straight.Move("a1", "b2") {
        t.Error("Should not move from a1 to b2")
    }

    straight.squares = 2
    if straight.Move("a1", "a7") {
        t.Error("Should not move from a1 to a7")
    }
    if straight.Move("a1", "f1") {
        t.Error("Should not move from a1 to f1")
    }
}

func TestDiagonalMovement(t *testing.T) {
    diagonal := Diagonal{}
    if !diagonal.Move("a1", "b2") {
        t.Error("Should move from a1 to b2")
    }
    if !diagonal.Move("a1", "h8") {
        t.Error("Should move from a1 to h8")
    }
    if !diagonal.Move("h8", "a1") {
        t.Error("Should move from h8 to a1")
    }
    if !diagonal.Move("a8", "h1") {
        t.Error("Should move from a8 to h1")
    }
    if diagonal.Move("a1", "b1") {
        t.Error("Should not move from a1 to b1")
    }
    if diagonal.Move("a1", "c1") {
        t.Error("Should not move from a1 to c1")
    }
    if diagonal.Move("a1", "a8") {
        t.Error("Should not move from a1 to a8")
    }
    if diagonal.Move("b1", "c3") {
        t.Error("Should not move from b1 to c3")
    }
    if diagonal.Move("a1", "b3") {
        t.Error("Should not move from a1 to b3")
    }

    diagonal.squares = 1
    if !diagonal.Move("a1", "b2") {
        t.Error("Should move from a1 to b2")
    }
    if !diagonal.Move("h8", "g7") {
        t.Error("Should move from h8 to g7")
    }
    if diagonal.Move("a1", "c3") {
        t.Error("Should not move from a1 to c3")
    }
    if diagonal.Move("c3", "a1") {
        t.Error("Should not move from c3 to a1")
    }
}

func TestLMovement(t *testing.T) {
    l := LMovement{}
    if !l.Move("b1", "c3") {
        t.Error("Should move from b1 to c3")
    }
    if !l.Move("b1", "a3") {
        t.Error("Should move from b1 to a3")
    }
    if !l.Move("c3", "b1") {
        t.Error("Should move from c3 to b1")
    }
    if l.Move("b1", "c2") {
        t.Error("Should not move from b1 to c2")
    }
}

func TestForwardMovement(t *testing.T) {
    forward := Forward{1}

    if !forward.Move("a1", "a2") {
        t.Error("Should move from a1 to a2")
    }
    if !forward.Move("a2", "a3") {
        t.Error("Should move from a2 to a3")
    }
    if !forward.Move("h1", "h2") {
        t.Error("Should move from h1 to h2")
    }
    if forward.Move("a1", "b1") {
        t.Error("Should not move from a1 to b1")
    }
    if forward.Move("a1", "a3") {
        t.Error("Should not move from a1 to a3")
    }
    if forward.Move("a2", "a1") {
        t.Error("Should not move from a2 to a1")
    }
    if forward.Move("a2", "b3") {
        t.Error("Should not move from a2 to b3")
    }
}
