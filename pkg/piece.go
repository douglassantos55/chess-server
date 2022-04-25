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

type AllowedMove struct {
	From Square
	To   Square
}

type Movement interface {
	HasMoves(from string, board *Board) bool
	IsValid(from, to string) bool
	CanCapture(from, to string, color Color, board *Board) bool
	IsAllowed(from, to string, color Color, board *Board) []AllowedMove
}

type Castle struct {
	origin Square
}

func (c Castle) CanCapture(from, to string, color Color, board *Board) bool {
	return false
}

func (c Castle) HasMoves(from string, board *Board) bool {
	return false
}

func (c Castle) IsValid(from, to string) bool {
	source, _ := parseSquare(from)
	dest, _ := parseSquare(to)

	colDistance := Abs(int(source.col - dest.col))
	rowDistance := source.row - dest.row

	return from == c.origin.String() && rowDistance == 0 && colDistance == 2
}

func (c Castle) IsAllowed(from, to string, color Color, board *Board) []AllowedMove {
	source, _ := parseSquare(from)
	dest, _ := parseSquare(to)

	if dest.col == 'g' && dest.row == c.origin.row {
		for square := source.Right(); square.col <= 'h'; square = square.Right() {
			if len(board.IsThreatened(square.String(), color)) > 0 {
				return []AllowedMove{}
			}

			piece := board.Square(square.String())
			if piece != Empty() {
				if square.col != 'h' || piece != Rook(color) {
					return []AllowedMove{}
				}
			}
		}

		return []AllowedMove{
			{From: source, To: dest},
			{From: Square{col: 'h', row: c.origin.row}, To: Square{col: 'f', row: c.origin.row}},
		}
	} else if dest.col == 'c' && dest.row == c.origin.row {
		for square := source.Left(); square.col >= 'a'; square = square.Left() {
			if len(board.IsThreatened(square.String(), color)) > 0 {
				return []AllowedMove{}
			}

			piece := board.Square(square.String())
			if piece != Empty() {
				if square.col != 'a' || piece != Rook(color) {
					return []AllowedMove{}
				}
			}
		}

		return []AllowedMove{
			{From: source, To: dest},
			{From: Square{col: 'a', row: c.origin.row}, To: Square{col: 'd', row: c.origin.row}},
		}
	}

	return []AllowedMove{}
}

type Forward struct {
	squares int
	moved   bool
}

func (f Forward) CanCapture(from, to string, color Color, board *Board) bool {
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr == nil || sourceErr == nil {
		piece := board.Square(to)
		colDistance := int(dest.col - source.col)

		return Abs(colDistance) == Abs(f.squares) && (piece != Empty() && piece.Color != color)
	}

	return false
}

func (f Forward) HasMoves(from string, board *Board) bool {
	piece := board.Square(from)
	square, err := parseSquare(from)

	if err != nil {
		return false
	}

	to := square.Up().String()

	return f.IsValid(from, to) && len(f.IsAllowed(from, to, piece.Color, board)) > 0 &&
		len(board.IsThreatened(to, piece.Color)) == 0
}

func (f Forward) IsAllowed(from, to string, color Color, board *Board) []AllowedMove {
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr == nil || sourceErr == nil {
		piece := board.Square(to)
		colDistance := int(dest.col - source.col)

		forward := colDistance == 0 && piece == Empty()
		capture := f.CanCapture(from, to, color, board)

		if forward || capture {
			// TODO: this doesnt work since it's no reference
			f.moved = true

			return []AllowedMove{{
				From: source,
				To:   dest,
			}}
		}
	}

	return []AllowedMove{}
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

func (s Straight) CanCapture(from, to string, color Color, board *Board) bool {
	return len(s.IsAllowed(from, to, color, board)) > 0
}

func (s Straight) HasMoves(from string, board *Board) bool {
	piece := board.Square(from)
	square, err := parseSquare(from)

	if piece.king {
		board.matrix[square.row][square.col] = Empty()
	}

	defer func() {
		if piece.king {
			board.matrix[square.row][square.col] = piece
		}
	}()

	if err != nil {
		return false
	}

	squares := []string{
		square.Up().String(),
		square.Down().String(),
		square.Left().String(),
		square.Right().String(),
	}

	for _, square := range squares {
		if len(s.IsAllowed(from, square, piece.Color, board)) > 0 && (!piece.king || len(board.IsThreatened(square, piece.Color)) == 0) {
			return true
		}
	}
	return false
}

func (s Straight) IsAllowed(from, to string, color Color, board *Board) []AllowedMove {
	moveRange, err := NewRange(from, to)

	if err != nil {
		return []AllowedMove{}
	}

	for !moveRange.Done() {
		cur := moveRange.Next()
		piece := board.matrix[cur.row][cur.col]

		if piece != Empty() {
			if piece.Color == color || cur != moveRange.until {
				return []AllowedMove{}
			}
		}
	}

	return []AllowedMove{{
		From: moveRange.from,
		To:   moveRange.until,
	}}
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

func (d Diagonal) CanCapture(from, to string, color Color, board *Board) bool {
	return len(d.IsAllowed(from, to, color, board)) > 0
}

func (d Diagonal) HasMoves(from string, board *Board) bool {
	piece := board.Square(from)
	source, err := parseSquare(from)

	if err != nil {
		return false
	}

	if piece.king {
		board.matrix[source.row][source.col] = Empty()
	}

	defer func() {
		if piece.king {
			board.matrix[source.row][source.col] = piece
		}
	}()

	squares := []string{
		source.UpRight().String(),
		source.DownRight().String(),
		source.UpLeft().String(),
		source.DownLeft().String(),
	}

	for _, square := range squares {
		if len(d.IsAllowed(from, square, piece.Color, board)) > 0 && (!piece.king || len(board.IsThreatened(square, piece.Color)) == 0) {
			return true
		}
	}
	return false
}

func (d Diagonal) IsAllowed(from, to string, color Color, board *Board) []AllowedMove {
	moveRange, err := NewRange(from, to)

	if err != nil {
		return []AllowedMove{}
	}

	for !moveRange.Done() {
		cur := moveRange.Next()
		piece := board.matrix[cur.row][cur.col]

		if piece != Empty() {
			if piece.Color == color || cur != moveRange.until {
				return []AllowedMove{}
			}
		}
	}

	return []AllowedMove{{
		From: moveRange.from,
		To:   moveRange.until,
	}}
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

func (l LMovement) CanCapture(from, to string, color Color, board *Board) bool {
	return len(l.IsAllowed(from, to, color, board)) > 0
}

func (l LMovement) HasMoves(from string, board *Board) bool {
	source, err := parseSquare(from)

	if err != nil {
		return false
	}

	piece := board.Square(from)

	squares := []string{
		fmt.Sprintf("%s%d", string(rune(source.col+1)), source.row+2),
		fmt.Sprintf("%s%d", string(rune(source.col-1)), source.row-2),
		fmt.Sprintf("%s%d", string(rune(source.col+1)), source.row-2),
		fmt.Sprintf("%s%d", string(rune(source.col-1)), source.row+2),
	}

	for _, square := range squares {
		if len(l.IsAllowed(from, square, piece.Color, board)) > 0 {
			return true
		}
	}
	return false
}

func (l LMovement) IsAllowed(from, to string, color Color, board *Board) []AllowedMove {
	piece := board.Square(to)
	if piece == Empty() || piece.Color != color {
		source, _ := parseSquare(from)
		dest, _ := parseSquare(to)

		return []AllowedMove{{
			From: source,
			To:   dest,
		}}
	}
	return []AllowedMove{}
}

func (l LMovement) IsValid(from, to string) bool {
	dest, destErr := parseSquare(to)
	source, sourceErr := parseSquare(from)

	if destErr != nil || sourceErr != nil {
		return false
	}

	rowDistance := Abs(dest.row - source.row)
	colDistance := Abs(int(dest.col - source.col))

	return (rowDistance == 2 && colDistance == 1) || (colDistance == 2 && rowDistance == 1)
}

type Combined struct {
	movements []Movement
}

func (c Combined) CanCapture(from, to string, color Color, board *Board) bool {
	for _, movement := range c.movements {
		if movement.CanCapture(from, to, color, board) {
			return true
		}
	}
	return false
}

func (c Combined) HasMoves(from string, board *Board) bool {
	for _, movement := range c.movements {
		if movement.HasMoves(from, board) {
			return true
		}
	}

	return false
}

func (c Combined) IsAllowed(from, to string, color Color, board *Board) []AllowedMove {
	result := []AllowedMove{}
	found := map[AllowedMove]bool{}

	for _, movement := range c.movements {
		moves := movement.IsAllowed(from, to, color, board)
		for _, move := range moves {
			if !found[move] {
				found[move] = true
				result = append(result, move)
			}
		}
	}
	return result
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

func (p *Piece) Move(from, to string, board *Board) []AllowedMove {
	if p.Movement.IsValid(from, to) {
		moves := p.Movement.IsAllowed(from, to, p.Color, board)
		isAllowed := (!p.king || len(board.IsThreatened(to, p.Color)) == 0)

		if isAllowed {
			return moves
		}
	}
	return []AllowedMove{}
}

func (p *Piece) CanCapture(from, to string, board *Board) bool {
	return p.Movement.IsValid(from, to) && p.Movement.CanCapture(from, to, p.Color, board)
}

func (p *Piece) CanMove(from, square string, board *Board) bool {
	return p.Movement.IsValid(from, square) && len(p.Movement.IsAllowed(from, square, p.Color, board)) > 0
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
	square := "e1"
	if color == Black {
		square = "e8"
	}

	origin, _ := parseSquare(square)

	movement := Combined{
		[]Movement{Straight{1}, Castle{origin}, Diagonal{1}},
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
