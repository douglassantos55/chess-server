package main

import (
	"math"
)

func Abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}

type Movement interface {
	IsValid(from, to string) bool
	IsAllowed(from, to string, color Color, board *Board) bool
}

type Forward struct {
	squares int
	moved   bool
}

func (f Forward) IsAllowed(from, to string, color Color, board *Board) bool {
	if board.Square(to) == Empty() {
		f.moved = true
	}
	return f.moved
}

func (f Forward) IsValid(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := int(destRow - fromRow)
	colDistance := Abs(int(destCol) - int(fromCol))

	return colDistance == 0 && (rowDistance == f.squares || (!f.moved && rowDistance == f.squares*2))
}

type Straight struct {
	squares int
}

func (s Straight) IsAllowed(from, to string, color Color, board *Board) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	startRow := math.Min(float64(destRow), float64(fromRow))
	endRow := math.Max(float64(destRow), float64(fromRow))

	startCol := math.Min(float64(destCol), float64(fromCol))
	endCol := math.Max(float64(destCol), float64(fromCol))

	for i := int(startRow); i <= int(endRow); i++ {
		for j := int(startCol); j <= int(endCol); j++ {
			piece := board.matrix[i][rune(j)]

			if (i != fromRow || j != int(fromCol)) && piece != Empty() {
				if piece.Color == color || i < int(endRow) || j < int(endCol) {
					return false
				}
			}
		}
	}

	return true
}

func (s Straight) IsValid(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := Abs(int(destRow - fromRow))
	colDistance := Abs(int(destCol) - int(fromCol))

	return (destRow == fromRow || destCol == fromCol) && (s.squares == 0 || (rowDistance <= s.squares && colDistance <= s.squares))
}

type Diagonal struct {
	squares int
}

func (d Diagonal) IsAllowed(from, to string, color Color, board *Board) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	startRow := math.Min(float64(destRow), float64(fromRow))
	startCol := math.Min(float64(destCol), float64(fromCol))

	distance := Abs(destRow - fromRow)

	for i := 0; i <= distance; i++ {
		piece := board.matrix[int(startRow)+i][rune(int(startCol)+i)]

		if piece != Empty() {
			if (piece.Color != color && i != distance) || (piece.Color == color && int(startRow)+i != fromRow) {
				return false
			}
		}
	}

	return true
}

func (d Diagonal) IsValid(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := Abs(destRow - fromRow)
	colDistance := Abs(int(destCol) - int(fromCol))

	return rowDistance == colDistance && (d.squares == 0 || (rowDistance <= d.squares && colDistance <= d.squares))
}

type LMovement struct{}

func (l LMovement) IsAllowed(from, to string, color Color, board *Board) bool {
	piece := board.Square(to)
	return piece == Empty() || piece.Color != color
}

func (l LMovement) IsValid(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := Abs(destRow - fromRow)
	colDistance := Abs(int(destCol) - int(fromCol))

	return rowDistance == 2 && colDistance == 1
}

type Combined struct {
	movements []Movement
}

func (c Combined) IsAllowed(from, to string, color Color, board *Board) bool {
	for _, movement := range c.movements {
		if movement.IsAllowed(from, to, color, board) {
			return true
		}
	}

	return false
}

func (c Combined) IsValid(from, to string) bool {
	for _, movement := range c.movements {
		if movement.IsValid(from, to) {
			return true
		}
	}
	return false
}

type Piece struct {
	Color    Color
	Notation string
	Movement Movement
	king     bool
}

func (p *Piece) Move(from, to string, board *Board) bool {
	return p.Sees(from, to, board) && (!p.king || !board.IsThreatned(to, p.Color))
}

func (p *Piece) Sees(from, square string, board *Board) bool {
	return p.Movement.IsValid(from, square) && p.Movement.IsAllowed(from, square, p.Color, board)
}

func CreatePiece(name string, color Color, movement Movement) Piece {
	return Piece{
		Color:    color,
		Notation: name,
		Movement: movement,
		king:     name == "K",
	}
}

func Empty() Piece {
	return CreatePiece("", White, nil)
}
func Rook(color Color) Piece {
	return CreatePiece("R", color, Straight{})
}
func Knight(color Color) Piece {
	return CreatePiece("N", color, LMovement{})
}
func Bishop(color Color) Piece {
	return CreatePiece("B", color, Diagonal{})
}
func Queen(color Color) Piece {
	movement := Combined{
		[]Movement{Straight{}, Diagonal{}},
	}
	return CreatePiece("Q", color, movement)
}
func King(color Color) Piece {
	movement := Combined{
		[]Movement{Straight{1}, Diagonal{1}},
	}
	return CreatePiece("K", color, movement)
}
func Pawn(color Color) Piece {
	direction := 1

	if color == Black {
		direction = -1
	}

	return CreatePiece("p", color, Forward{squares: direction})
}
