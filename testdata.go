package mmr

// Pos, Height, HasChildren, Left, Right, Parent
var Sequence = []Node{
	{0, 0, 2, 0, 0},
	{1, 0, 2, 0, 0},
	{2, 1, 6, 0, 1},
	{3, 0, 5, 0, 0},
	{4, 0, 5, 0, 0},
	{5, 1, 6, 3, 4},
	{6, 2, 14, 2, 5},
	{7, 0, 9, 0, 0},
	{8, 0, 9, 0, 0},
	{9, 1, 13, 7, 8},
	{10, 0, 12, 0, 0},
	{11, 0, 12, 0, 0},
	{12, 1, 13, 10, 11},
	{13, 2, 14, 9, 12},
	{14, 3, 30, 6, 13},
	{15, 0, 17, 0, 0},
	{16, 0, 17, 0, 0},
	{17, 1, 21, 15, 16},
}
