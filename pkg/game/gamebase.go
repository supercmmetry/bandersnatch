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
	Id              uint64    `json:"id"`
	Data            *NodeData `json:"data"`
	Leader          *Node
	LeftChild       *Node
	RightChild      *Node
	LeftNodeIds     []uint64 `json:"left_nodes"`
	RightNodeIds    []uint64 `json:"right_nodes"`
	IsLeader        bool     `json:"is_leader"`
	IsLeaf          bool     `json:"is_leaf"`
	CanHoldArtifact bool     `json:"can_hold_artifact"`
	RandomizePath   bool     `json:"randomize_path"`
}

func (n *Node) Traverse(opt Option) *Node {
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

func (n *Node) GetNodeByNum(num int) *Node {
	curr := n
	queue := make([]*Node, 0)
	queue = append(queue, curr)
	visited := make(map[uint64]struct{})
	for len(queue) > 0 && num >= 0 {
		curr = queue[0]
		queue = queue[1:]
		if _, ok := visited[curr.Id]; ok {
			continue
		}
		visited[curr.Id] = struct{}{}

		num--
		if curr.LeftChild != nil {
			queue = append(queue, curr.LeftChild)
		}
		if curr.RightChild != nil {
			queue = append(queue, curr.RightChild)
		}
	}
	return curr
}

type Artifact struct {
	Id                  uint64  `json:"id"`
	ScrambleCoefficient float64 `json:"scramble_coeff"`
	Description         string  `json:"description"`
}
