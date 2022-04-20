package pkg

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
	piece := board.Square(from)
	square, err := parseSquare(from)

	if err != nil {
		return false
	}

	to := square.Up()
	return f.IsAllowed(from, to.String(), piece.Color, board) && len(board.IsThreatened(to.String(), piece.Color)) == 0
}

func (f Forward) IsAllowed(from, to string, color Color, board *Board) bool {
	piece := board.Square(to)
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr != nil || sourceErr != nil {
		return false
	}

	colDistance := int(dest.col - source.col)
	forward := colDistance == 0 && piece == Empty()
	capture := Abs(colDistance) == Abs(f.squares) && (piece != Empty() && piece.Color != color)

	if forward || capture {
		f.moved = true
	}

	return f.moved
}

func (f Forward) IsValid(from, to string) bool {
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr != nil || sourceErr != nil {
		return false
	}

	rowDistance := int(dest.row - source.row)
	return rowDistance == f.squares || (!f.moved && rowDistance == f.squares*2)
}

type Straight struct {
	squares int
}

func (s Straight) HasMoves(from string, board *Board) bool {
	square, err := parseSquare(from)

	if err != nil {
		return false
	}

	piece := board.Square(from)

	up := square.Up().String()
	down := square.Down().String()
	left := square.Left().String()
	right := square.Right().String()

	return ((s.IsAllowed(from, up, piece.Color, board) && len(board.IsThreatened(up, piece.Color)) == 0) ||
		(s.IsAllowed(from, down, piece.Color, board) && len(board.IsThreatened(down, piece.Color)) == 0) ||
		(s.IsAllowed(from, left, piece.Color, board) && len(board.IsThreatened(left, piece.Color)) == 0) ||
		(s.IsAllowed(from, right, piece.Color, board) && len(board.IsThreatened(right, piece.Color)) == 0))

}

func (s Straight) IsAllowed(from, to string, color Color, board *Board) bool {
	moveRange, err := NewRange(from, to)

	if err != nil {
		return false
	}

	for !moveRange.Done() {
		cur := moveRange.Next()
		piece := board.matrix[cur.row][cur.col]

		if piece != Empty() {
			if piece.Color == color || cur != moveRange.until {
				return false
			}
		}
	}

	return true
}

func (s Straight) IsValid(from, to string) bool {
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr != nil || sourceErr != nil {
		return false
	}

	rowDistance := Abs(dest.row - source.row)
	colDistance := Abs(int(dest.col - source.col))

	return (dest.row == source.row || dest.col == source.col) &&
		(s.squares == 0 || (rowDistance <= s.squares && colDistance <= s.squares))
}

type Diagonal struct {
	squares int
}

func (d Diagonal) HasMoves(from string, board *Board) bool {
	source, err := parseSquare(from)

	if err != nil {
		return false
	}

	piece := board.Square(from)

	upRight := source.UpRight()
	downRight := source.DownRight()
	upLeft := source.UpLeft()
	downLeft := source.DownLeft()

	return ((d.IsAllowed(from, upRight.String(), piece.Color, board) && len(board.IsThreatened(upRight.String(), piece.Color)) == 0) ||
		(d.IsAllowed(from, downRight.String(), piece.Color, board) && len(board.IsThreatened(downRight.String(), piece.Color)) == 0) ||
		(d.IsAllowed(from, upLeft.String(), piece.Color, board) && len(board.IsThreatened(upLeft.String(), piece.Color)) == 0) ||
		(d.IsAllowed(from, downLeft.String(), piece.Color, board) && len(board.IsThreatened(downLeft.String(), piece.Color)) == 0))
}

func (d Diagonal) IsAllowed(from, to string, color Color, board *Board) bool {
	moveRange, err := NewRange(from, to)

	if err != nil {
		return false
	}

	for !moveRange.Done() {
		cur := moveRange.Next()
		piece := board.matrix[cur.row][cur.col]

		if piece != Empty() {
			if piece.Color == color || cur != moveRange.until {
				return false
			}
		}
	}

	return true
}

func (d Diagonal) IsValid(from, to string) bool {
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr != nil || sourceErr != nil {
		return false
	}

	rowDistance := Abs(dest.row - source.row)
	colDistance := Abs(int(dest.col - source.col))

	return rowDistance == colDistance && (d.squares == 0 || (rowDistance <= d.squares && colDistance <= d.squares))
}

type LMovement struct{}

func (l LMovement) HasMoves(from string, board *Board) bool {
	source, err := parseSquare(from)

	if err != nil {
		return false
	}

	piece := board.Square(from)

	up := fmt.Sprintf("%s%d", string(rune(source.col+1)), source.row+2)
	down := fmt.Sprintf("%s%d", string(rune(source.col-1)), source.row-2)
	left := fmt.Sprintf("%s%d", string(rune(source.col+1)), source.row-2)
	right := fmt.Sprintf("%s%d", string(rune(source.col-1)), source.row+2)

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
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr != nil || sourceErr != nil {
		return false
	}

	rowDistance := Abs(dest.row - source.row)
	colDistance := Abs(int(dest.col - source.col))

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
	sees := p.Sees(from, to, board)
	isAllowed := (!p.king || len(board.IsThreatened(to, p.Color)) == 0)

	return sees && isAllowed
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
