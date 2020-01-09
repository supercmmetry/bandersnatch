package core

type NodeState struct {
	Tag              string
	ContainsArtifact bool
	LeapToNexus      bool
	LeapToNode       *Node
	JumpBack         bool
	TransitionText   string
}

type NodeDisplay struct {
	Text        string
	OptionLeft  string
	OptionRight string
}

type Node struct {
	Display    NodeDisplay
	State      NodeState
	LeftChild  *Node
	RightChild *Node
	Parent     *Node
	IsLeaf     bool
	IsLeader   bool
}

const (
	BandersnatchLeft        = true
	BandersnatchRight       = false
	BandersnatchRewindLimit = 5
)
