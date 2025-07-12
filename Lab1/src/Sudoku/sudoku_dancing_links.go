package main

type Node struct {
	Left  *Node
	Right *Node
	Up    *Node
	Down  *Node
	Col   *Node //Column
	Name  int
	Size  int
}

type Column = Node

const (
	ROW         = 9
	COL         = 9
	N           = 81
	NEIGHBOR    = 20
	SOLVEPTHNUM = 20
	KMaxNodes   = 1 + 81*4 + 9*9*9*4
	KMaxColumns = 400
	KRow        = 100
	KCol        = 200
	KBox        = 300
)

var Board [SOLVEPTHNUM][N]int

type Dance struct {
	Root_    *Column
	Inout_   []int
	Columns_ [KMaxColumns]*Column
	Stack_   []*Node
	Nodes_   [KMaxNodes]Node
	CurNode_ int
}

func (d *Dance) NewColumn(n int) *Column {
	// assert.Assert(d.CurNode_ < KMaxNodes)
	c := &d.Nodes_[d.CurNode_]
	d.CurNode_++
	c.Left = c
	c.Right = c
	c.Up = c
	c.Down = c
	c.Col = c
	c.Name = n
	return c
}

func (d *Dance) AppendColumn(n int) {
	// assert.Assert(d.Columns_[n] == nil)

	c := d.NewColumn(n)
	d.PutLeft(d.Root_, c)
	d.Columns_[n] = c
}

func (d *Dance) NewRow(col int) *Node {
	// assert.Assert(d.Columns_[col] != nil)
	// assert.Assert(d.CurNode_ < KMaxNodes)

	r := &d.Nodes_[d.CurNode_]
	d.CurNode_++

	r.Left = r
	r.Right = r
	r.Up = r
	r.Down = r
	r.Name = col
	r.Col = d.Columns_[col]
	d.PutUp(r.Col, r)
	return r
}

func (d *Dance) GetRowCol(row, val int) int {
	return KRow + row*10 + val
}

func (d *Dance) GetColCol(col, val int) int {
	return KCol + col*10 + val
}

func (d *Dance) GetBoxCol(box, val int) int {
	return KBox + box*10 + val
}
func NewDance(inout []int) *Dance {
	d := &Dance{
		Inout_:   inout,
		CurNode_: 0,
	}

	d.Root_ = d.NewColumn(0)
	d.Root_.Left = d.Root_
	d.Root_.Right = d.Root_

	for i := 0; i < KMaxColumns; i++ {
		d.Columns_[i] = nil
	}

	rows := make([][10]bool, ROW)
	cols := make([][10]bool, COL)
	boxes := make([][10]bool, ROW)

	for i := 0; i < N; i++ {
		row := i / 9
		col := i % 9
		box := row/3*3 + col/3
		val := inout[i]
		rows[row][val] = true
		cols[col][val] = true
		boxes[box][val] = true
	}

	for i := 0; i < N; i++ {
		if inout[i] == 0 {
			d.AppendColumn(i)
		}
	}

	for i := 0; i < 9; i++ {
		for v := 1; v < 10; v++ {
			if !rows[i][v] {
				d.AppendColumn(d.GetRowCol(i, v))
			}
			if !cols[i][v] {
				d.AppendColumn(d.GetColCol(i, v))
			}
			if !boxes[i][v] {
				d.AppendColumn(d.GetBoxCol(i, v))
			}
		}
	}

	for i := 0; i < N; i++ {
		if inout[i] == 0 {
			row := i / 9
			col := i % 9
			box := row/3*3 + col/3

			for v := 1; v < 10; v++ {
				if !(rows[row][v] || cols[col][v] || boxes[box][v]) {
					n0 := d.NewRow(i)
					nr := d.NewRow(d.GetRowCol(row, v))
					nc := d.NewRow(d.GetColCol(col, v))
					nb := d.NewRow(d.GetBoxCol(box, v))
					d.PutLeft(n0, nr)
					d.PutLeft(n0, nc)
					d.PutLeft(n0, nb)
				}
			}
		}
	}

	return d
}
func (d *Dance) GetMinColumn() *Column {
	c := d.Root_.Right
	minSize := c.Size
	if minSize > 1 {
		for cc := c.Right; cc != d.Root_; cc = cc.Right {
			if minSize > cc.Size {
				c = cc
				minSize = cc.Size
				if minSize <= 1 {
					break
				}
			}
		}
	}
	return c
}

func (d *Dance) Cover(c *Column) {
	c.Right.Left = c.Left
	c.Left.Right = c.Right

	for row := c.Down; row != c; row = row.Down {
		for j := row.Right; j != row; j = j.Right {
			j.Down.Up = j.Up
			j.Up.Down = j.Down
			j.Col.Size--
		}
	}
}

func (d *Dance) Uncover(c *Column) {
	for row := c.Up; row != c; row = row.Up {
		for j := row.Left; j != row; j = j.Left {
			j.Col.Size++
			j.Down.Up = j
			j.Up.Down = j
		}
	}
	c.Right.Left = c
	c.Left.Right = c
}

func (d *Dance) Solve() bool {
	if d.Root_.Left == d.Root_ {
		for i := 0; i < len(d.Stack_); i++ {
			n := d.Stack_[i]
			cell := -1
			val := -1
			for cell == -1 || val == -1 {
				if n.Name < 100 {
					cell = n.Name
				} else {
					val = n.Name % 10
				}
				n = n.Right
			}
			d.Inout_[cell] = val
		}
		return true
	}

	col := d.GetMinColumn()
	d.Cover(col)
	for row := col.Down; row != col; row = row.Down {
		d.Stack_ = append(d.Stack_, row)
		for j := row.Right; j != row; j = j.Right {
			d.Cover(j.Col)
		}
		if d.Solve() {
			return true
		}
		d.Stack_ = d.Stack_[:len(d.Stack_)-1]
		for j := row.Left; j != row; j = j.Left {
			d.Uncover(j.Col)
		}
	}
	d.Uncover(col)
	return false
}

func (d *Dance) PutLeft(old, nnew *Column) {
	nnew.Left = old.Left
	nnew.Right = old
	old.Left.Right = nnew
	old.Left = nnew
}

func (d *Dance) PutUp(old *Column, nnew *Node) {
	nnew.Up = old.Up
	nnew.Down = old
	old.Up.Down = nnew
	old.Up = nnew
	old.Size++
	nnew.Col = old
}

func SolveSudokuDancingLinks(unused int, threadId int) bool {
	// mutex.Lock()
	// defer mutex.Unlock()
	// board := Board[threadId]
	// fmt.Println(board)
	d := NewDance(Board[threadId][:])
	// fmt.Println(board)
	return d.Solve()
}
