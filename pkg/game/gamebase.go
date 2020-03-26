package game

type Option uint

const (
	OptionLeft Option = iota
	OptionRight
)

type NodeData struct {
	Question    string `json:"question"`
	LeftOption  string `json:"left_option"`
	RightOption string `json:"right_option"`
}

type Node struct {
	Id                  uint64    `json:"id"`
	Data                *NodeData `json:"data"`
	Leader              *Node
	LeftChild           *Node
	RightChild          *Node
	LeftNodeIds         []uint64 `json:"left_nodes"`
	RightNodeIds        []uint64 `json:"right_nodes"`
	ArtifactIds         []uint64 `json:"artifact_ids"`
	RequiredArtifactIds []uint64 `json:"required_artifact_ids"`
	IsLeader            bool     `json:"is_leader"`
	IsLeaf              bool     `json:"is_leaf"`
	RandomizePath       bool
}

func (n Node) Traverse(opt Option) *Node {
	if opt == OptionLeft {
		if n.LeftChild == nil {
			return n.RightChild
		}
		return n.LeftChild
	} else {
		if n.RightChild == nil {
			return n.LeftChild
		}
		return n.RightChild
	}
}

type Artifact struct {
	Id           uint64   `json:"id"`
	Description  string   `json:"description"`
	Score        uint64   `json:"score"`
}
