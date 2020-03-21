package bandersnatch

type Option uint
const (
	OptionLeft Option = iota
	OptionRight
)

type NodeData struct {
	Question    string
	LeftOption  string
	RightOption string
}

type Node struct {
	Id         uint64
	Data       *NodeData
	Parent     *Node
	LeftChild  *Node
	RightChild *Node
	IsLeader   bool
	IsLeaf     bool
}

func (n *Node) FetchLeader() *Node {
	curr := n
	for curr.Parent != nil {
		curr = curr.Parent
	}
	return curr
}

func (n *Node) Traverse(opt Option) {
	if opt == OptionLeft {
		*n = *n.LeftChild
	} else {
		*n = *n.RightChild
	}
}

func (n *Node) GetNodeByNum(num int) *Node {
	newNode := n
	for num != 0 && !newNode.IsLeaf{
		lsb := num & 1
		if lsb == 1 {
			newNode = newNode.RightChild
		} else {
			newNode = newNode.LeftChild
		}

		num >>= 1
	}

	return newNode
}

type Artifact struct {
	Id                  uint64
	ScrambleCoefficient float64
	Description         string
}
