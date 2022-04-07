package main

import (
	"fmt"
)

func Abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}

type Movement interface {
	HasMoves(from string, board *Board) bool
	IsValid(from, to string) bool
	IsAllowed(from, to string, color Color, board *Board) bool
}

type Forward struct {
	squares int
	moved   bool
}

func (f Forward) HasMoves(from string, board *Board) bool {
	row, col := parseSquare(from)
	piece := board.Square(from)
	to := fmt.Sprintf("%s%d", string(rune(col)), row+1)

	return f.IsAllowed(from, to, piece.Color, board) && !board.IsThreatned(to, piece.Color)
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

func (s Straight) HasMoves(from string, board *Board) bool {
	row, col := parseSquare(from)
	piece := board.Square(from)

	up := fmt.Sprintf("%s%d", string(rune(col)), row+2)
	down := fmt.Sprintf("%s%d", string(rune(col)), row)
	left := fmt.Sprintf("%s%d", string(rune(col-1)), row+1)
	right := fmt.Sprintf("%s%d", string(rune(col+1)), row+1)

	return ((s.IsAllowed(from, up, piece.Color, board) && !board.IsThreatned(up, piece.Color)) ||
		(s.IsAllowed(from, down, piece.Color, board) && !board.IsThreatned(down, piece.Color)) ||
		(s.IsAllowed(from, left, piece.Color, board) && !board.IsThreatned(left, piece.Color)) ||
		(s.IsAllowed(from, right, piece.Color, board) && !board.IsThreatned(right, piece.Color)))

}

func (s Straight) IsAllowed(from, to string, color Color, board *Board) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	if destRow < 0 || destRow > 7 || destCol < 'a' || destCol > 'h' {
		return false
	}

	rowStep := 1
	colStep := 1

	if destRow < fromRow {
		rowStep = -1
	}

	if destCol < fromCol {
		colStep = -1
	}

	rowDistance := int(float64(Abs(destRow - fromRow)))
	colDistance := int(float64(Abs(int(destCol - fromCol))))

	for i := 0; i <= rowDistance; i++ {
		for j := 0; j <= colDistance; j++ {
			r := fromRow + rowStep*i
			c := int(fromCol) + colStep*j

			piece := board.matrix[r][rune(c)]

			if piece != Empty() {
				if (i != 0 || j != 0) && (piece.Color == color || (i != rowDistance || j != colDistance)) {
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

func (d Diagonal) HasMoves(from string, board *Board) bool {
	row, col := parseSquare(from)
	piece := board.Square(from)

	up := fmt.Sprintf("%s%d", string(rune(col+1)), row+2)
	down := fmt.Sprintf("%s%d", string(rune(col-1)), row)
	left := fmt.Sprintf("%s%d", string(rune(col-1)), row+2)
	right := fmt.Sprintf("%s%d", string(rune(col+1)), row)

	return ((d.IsAllowed(from, up, piece.Color, board) && !board.IsThreatned(up, piece.Color)) ||
		(d.IsAllowed(from, down, piece.Color, board) && !board.IsThreatned(down, piece.Color)) ||
		(d.IsAllowed(from, left, piece.Color, board) && !board.IsThreatned(left, piece.Color)) ||
		(d.IsAllowed(from, right, piece.Color, board) && !board.IsThreatned(right, piece.Color)))
}

func (d Diagonal) IsAllowed(from, to string, color Color, board *Board) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	if destRow < 0 || destRow > 7 || destCol < 'a' || destCol > 'h' {
		return false
	}

	rowStep := 1
	colStep := 1

	if destCol < fromCol {
		colStep = -1
	}
	if destRow < fromRow {
		rowStep = -1
	}

	distance := Abs(destRow - fromRow)

	for i := 1; i <= distance; i++ {
		piece := board.matrix[int(fromRow)+rowStep*i][rune(int(fromCol)+colStep*i)]
		if piece != Empty() {
			if piece.Color == color || i != distance {
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

func (l LMovement) HasMoves(from string, board *Board) bool {
	row, col := parseSquare(from)
	piece := board.Square(from)

	up := fmt.Sprintf("%s%d", string(rune(col+1)), row+2)
	down := fmt.Sprintf("%s%d", string(rune(col-1)), row-2)
	left := fmt.Sprintf("%s%d", string(rune(col+1)), row-2)
	right := fmt.Sprintf("%s%d", string(rune(col-1)), row+2)

	return (l.IsAllowed(from, up, piece.Color, board) ||
		l.IsAllowed(from, down, piece.Color, board) ||
		l.IsAllowed(from, left, piece.Color, board) ||
		l.IsAllowed(from, right, piece.Color, board))
}

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

func (c Combined) HasMoves(from string, board *Board) bool {
	for _, movement := range c.movements {
		if movement.HasMoves(from, board) {
			return true
		}
	}

	return false
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

func (p *Piece) HasMoves(from string, board *Board) bool {
	return p.Movement.HasMoves(from, board)
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
