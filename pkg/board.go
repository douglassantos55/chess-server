package pkg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func parseSquare(square string) (Square, error) {
	coord := strings.Split(square, "")
	row, _ := strconv.ParseInt(coord[1], 10, 64)
	col := []rune(coord[0])[0]

	if row < 1 || row > 8 || col < 'a' || col > 'h' {
		return Square{}, errors.New("Invalid square")
	}

	return Square{
		row: int(row - 1),
		col: col,
	}, nil
}

type Square struct {
	col rune
	row int
}

func (s Square) String() string {
	return fmt.Sprintf("%s%d", string(s.col), s.row+1)
}

func (s Square) Up() Square {
	return Square{
		col: s.col,
		row: s.row + 1,
	}
}
func (s Square) Down() Square {
	return Square{
		col: s.col,
		row: s.row - 1,
	}
}
func (s Square) Left() Square {
	return Square{
		col: s.col - 1,
		row: s.row,
	}
}
func (s Square) Right() Square {
	return Square{
		col: s.col + 1,
		row: s.row,
	}
}
func (s Square) UpRight() Square {
	return Square{
		col: s.col + 1,
		row: s.row + 1,
	}
}
func (s Square) DownRight() Square {
	return Square{
		col: s.col + 1,
		row: s.row - 1,
	}
}
func (s Square) UpLeft() Square {
	return Square{
		col: s.col - 1,
		row: s.row + 1,
	}
}
func (s Square) DownLeft() Square {
	return Square{
		col: s.col - 1,
		row: s.row - 1,
	}
}

type Range struct {
	from  Square
	until Square
	cur   Square
}

func NewRange(from, to string) (Range, error) {
	dest, err := parseSquare(to)
	if err != nil {
		return Range{}, err
	}

	source, err := parseSquare(from)
	if err != nil {
		return Range{}, err
	}

	return Range{
		until: dest,
		from:  source,
		cur:   source,
	}, nil
}

func (r Range) Done() bool {
	return r.cur == r.until
}

func (r *Range) Next() Square {
	if r.from.row > r.until.row {
		r.cur = r.cur.Down()
	} else if r.from.row < r.until.row {
		r.cur = r.cur.Up()
	}

	if r.from.col > r.until.col {
		r.cur = r.cur.Left()
	} else if r.from.col < r.until.col {
		r.cur = r.cur.Right()
	}

	return r.cur
}

type Board struct {
	mutex  *sync.Mutex
	matrix [8]map[rune]Piece
}

func NewBoard() *Board {
	return &Board{
		mutex: new(sync.Mutex),
		matrix: [8]map[rune]Piece{
			{
				'a': Rook(White),
				'b': Knight(White),
				'c': Bishop(White),
				'd': Queen(White),
				'e': King(White),
				'f': Bishop(White),
				'g': Knight(White),
				'h': Rook(White),
			},
			{
				'a': Pawn(White),
				'b': Pawn(White),
				'c': Pawn(White),
				'd': Pawn(White),
				'e': Pawn(White),
				'f': Pawn(White),
				'g': Pawn(White),
				'h': Pawn(White),
			},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{
				'a': Pawn(Black),
				'b': Pawn(Black),
				'c': Pawn(Black),
				'd': Pawn(Black),
				'e': Pawn(Black),
				'f': Pawn(Black),
				'g': Pawn(Black),
				'h': Pawn(Black),
			},
			{
				'a': Rook(Black),
				'b': Knight(Black),
				'c': Bishop(Black),
				'd': Queen(Black),
				'e': King(Black),
				'f': Bishop(Black),
				'g': Knight(Black),
				'h': Rook(Black),
			},
		},
	}
}

func (b *Board) Square(square string) Piece {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	s, _ := parseSquare(square)
	return b.matrix[s.row][s.col]
}

func (b *Board) CanBlock(threats []Range, color Color) bool {
	for _, threat := range threats {
		for !threat.Done() {
			for i := int('a'); i <= int('h'); i++ {
				for j := 1; j <= 8; j++ {
					piece := b.matrix[j-1][rune(i)]

					if piece != Empty() && !piece.king && piece.Color == color {
						from := fmt.Sprintf("%s%d", string(rune(i)), j)

						if piece.Sees(from, threat.cur.String(), b) {
							return true
						}
					}
				}
			}
			threat.Next()
		}
	}
	return false
}

func (b *Board) IsThreatened(square string, color Color) []Range {
	ranges := []Range{}

	for i := int('a'); i <= int('h'); i++ {
		for j := 1; j <= 8; j++ {
			piece := b.matrix[j-1][rune(i)]
			from := fmt.Sprintf("%s%d", string(rune(i)), j)

			if piece != Empty() && piece.Color != color {
				if piece.Sees(from, square, b) {
					threatRange, _ := NewRange(from, square)
					ranges = append(ranges, threatRange)
				}
			}
		}
	}

	return ranges
}

func (b *Board) Move(from, to string) {
	piece := b.Square(from)

	if piece.Move(from, to, b) {
		dest, _ := parseSquare(to)
		source, _ := parseSquare(from)

		b.matrix[dest.row][dest.col] = piece
		b.matrix[source.row][source.col] = Empty()
	}
}
