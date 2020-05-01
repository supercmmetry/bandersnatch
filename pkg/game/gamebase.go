package game

type Option uint

const (
	OptionLeft Option = iota
	OptionRight
)

type NodeData struct {
	Content       string                 `json:"content"`
	Miscellaneous map[string]interface{} `json:"misc"`
}

type Node struct {
	Id                       uint64    `json:"id"`
	Data                     *NodeData `json:"data"`
	Leader                   *Node
	LeftChild                *Node
	RightChild               *Node
	LeftNodeIds              []uint64 `json:"left_nodes"`
	RightNodeIds             []uint64 `json:"right_nodes"`
	ArtifactIds              []uint64 `json:"artifact_ids"`
	RequiredArtifactIds      []uint64 `json:"required_artifact_ids"`
	AntiRequisiteArtifactIds []uint64 `json:"anti_requisite_artifact_ids"`
	IsLeader                 bool     `json:"is_leader"`
	IsLeaf                   bool     `json:"is_leaf"`
	RandomizePath            bool
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
	Id            uint64                 `json:"id"`
	Miscellaneous map[string]interface{} `json:"misc"`
	Score         uint64                 `json:"score"`
}
