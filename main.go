package main

import (
	"fmt"
)

type Node struct {
	row      int
	reel     int
	rows     int
	reels    int
	parent   *Node
	children []*Node
	payline  []int
}

// TODO: Possible to skip storing parent and payline?

// A node represents a cell in the grid of rows/reels, and stores all
// the walk-paths that can start from that cell, in the form of child
// nodes who store their own child nodes and so on.
func NewNode(row, reel, rows, reels int, parent *Node) *Node {
	// Q:  any diff between Node{...} and  NewNode(...)?

	// A: I think fundamentally no.  But note that we are inside the
	// constructor NewNode so we cannot call NewNode yet.  The &Node{}
	// expression below is just a "node literal", like you have slice literals
	// []int{1, 2, 3}.  The thing is, we wouldn't need a NewNode constructor at
	// all if this was all we wanted to set when instantiating a node, but we
	// also want to do other things (seeing as there's more code below that
	// expression), thus we need a constructor function.
	node := &Node{
		row:      row,
		reel:     reel,
		rows:     rows,
		reels:    reels,
		parent:   parent,
		children: []*Node{},
	}
	if parent == nil {
		node.payline = []int{row}
	} else {
		node.payline = append(parent.payline, row)
	}
	if reel < reels {
		// Add 2 or 3 new children -- note the recursion, since this method
		// itself calls NewNode.  The method also takes the liberty of setting
		// our field `children`, which I feel could be more clear.

		// Q:  That's also why I have to use pointers/values *Node &Node, I
		// think.  If I write keyseqsfinder.go, find out if can avoid that.
		node.addChildren()
	}
	return node
}

// This is the slow part.  Nearly O(3^n), where n is the number of
// reels, because for each reel, the walk splits into three new paylines
// (sometimes just two).  3 -> 9 -> 27 -> ...
func (n *Node) addChildren() {
	minRow := n.row - 1
	if minRow < 0 {
		minRow = 0
	}
	maxRow := n.row + 1
	if maxRow >= n.rows {
		maxRow = n.rows - 1
	}
	for row := minRow; row <= maxRow; row++ {
		child := NewNode(row, n.reel+1, n.rows, n.reels, n)
		n.children = append(n.children, child)
	}
}

// Q:  is it necessary to use [][]int, couldn't it just be []int?
func (n *Node) getPaylines() [][]int {
	if len(n.children) == 0 {
		return [][]int{n.payline}
	}
	var paylines [][]int
	for _, child := range n.children {
		paylines = append(paylines, child.getPaylines()...)
	}
	return paylines
}

// I can probably get rid of Machine, it may have been an overly OOP way of
// doing things.
type Machine struct {
	rows  int
	reels int
	nodes []*Node
}

// here it does not need *Machine or &Machine -- but it's consistent with the
// other struct we have
func NewMachine(rows, reels int) *Machine {
	machine := &Machine{
		rows:  rows,
		reels: reels,
		nodes: []*Node{},
	}
	machine.addNodes()
	return machine
}

func (m *Machine) addNodes() {
	for row := 0; row < m.rows; row++ {
		node := NewNode(row, 1, m.rows, m.reels, nil)
		m.nodes = append(m.nodes, node)
	}
}

func (m *Machine) getPaylines() [][]int {
	var paylines [][]int
	for _, node := range m.nodes {
		paylines = append(paylines, node.getPaylines()...)
	}
	return paylines
}

// Used for performance profiling in main_test.go
func calcPaylines(rows, reels int) {
	_ = elkWays(rows, reels)
}

func elkWays(rows, reels int) [][]int {
	machine := NewMachine(rows, reels)
	return machine.getPaylines()
}

// For a nice printout from go run main.go
func main() {
	paylines := elkWays(4, 5)
	for _, p := range paylines {
		fmt.Println(p)
	}
	// Should be 178 paylines for 4 rows, 5 reels
	fmt.Println(len(paylines))
}
