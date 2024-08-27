package main

import (
	"fmt"
	"testing"
)

var table = []struct {
	rows  int
	reels int
}{
	{rows: 1, reels: 5},
	{rows: 10, reels: 5},
	{rows: 100, reels: 5},
	{rows: 1000, reels: 5},
	{rows: 10000, reels: 5},
	// {rows: 100000, reels: 5},
	// {rows: 1000000, reels: 5},

	// Reels is clearly the limiting factor. O(3^n) or so.
	{rows: 10, reels: 5},
	{rows: 10, reels: 10},
	// {rows: 10, reels: 15},
	// {rows: 10, reels: 20},
	// {rows: 10, reels: 25},
	// {rows: 10, reels: 30},
	// {rows: 10, reels: 35},
}

func BenchmarkPaylines(b *testing.B) {
	for _, t := range table {
		b.Run(fmt.Sprintf("rows %d reels %d", t.rows, t.reels),
			func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					calcPaylines(t.rows, t.reels)
				}
			})
	}
}

// To benchmark, run in shell:
// go test -bench=. -benchtime=2s *.go
