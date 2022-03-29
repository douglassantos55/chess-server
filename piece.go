package main

type Movement interface {
	Move(from, to string) bool
}

type Forward struct {
	squares int
}

func (f Forward) Move(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := int(destRow - fromRow)
	colDistance := Abs(int(destCol) - int(fromCol))

	return colDistance == 0 && rowDistance == f.squares
}

type Straight struct {
	squares int
}

func (s Straight) Move(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := Abs(int(destRow - fromRow))
	colDistance := Abs(int(destCol) - int(fromCol))

	return (destRow == fromRow || destCol == fromCol) && (s.squares == 0 || (rowDistance <= s.squares && colDistance <= s.squares))
}

type Diagonal struct {
	squares int
}

func (d Diagonal) Move(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := Abs(destRow - fromRow)
	colDistance := Abs(int(destCol) - int(fromCol))

	return rowDistance == colDistance && (d.squares == 0 || (rowDistance <= d.squares && colDistance <= d.squares))
}

type LMovement struct{}

func (l LMovement) Move(from, to string) bool {
	destRow, destCol := parseSquare(to)
	fromRow, fromCol := parseSquare(from)

	rowDistance := Abs(destRow - fromRow)
	colDistance := Abs(int(destCol) - int(fromCol))

	return rowDistance == 2 && colDistance == 1
}

type Piece struct {
	Color    Color
	Notation string
	Movement Movement
}

func (p *Piece) Move(from, to string) bool {
	return p.Movement.Move(from, to)
}

func CreatePiece(name string, color Color, movement Movement) Piece {
	return Piece{
		Color:    color,
		Notation: name,
		Movement: movement,
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
	return CreatePiece("Q", color, nil)
}
func King(color Color) Piece {
	return CreatePiece("K", color, nil)
}
func Pawn(color Color) Piece {
	return CreatePiece("p", color, Forward{1})
}
