package testdata

import "github.com/vsekhar/mmr"

// Pos, Height, HasChildren
var Sequence = []mmr.Node{
	{0, 0},
	{1, 0},
	{2, 1},
	{3, 0},
	{4, 0},
	{5, 1},
	{6, 2},
	{7, 0},
	{8, 0},
	{9, 1},
	{10, 0},
	{11, 0},
	{12, 1},
	{13, 2},
	{14, 3},
	{15, 0},
	{16, 0},
	{17, 1},
}

var FirstNode = mmr.Node{
	Pos:    0,
	Height: 0,
}
