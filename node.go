package mmr

type Node struct {
	Pos         int
	Height      int
	Parent      int
	HasChildren bool
	Left, Right int
}
