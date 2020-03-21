package bandersnatch

import (
	"encoding/json"
	"math/rand"
	"os"
)

type Player struct {
	Id                   uint64
	ArtifactDistribution map[*Artifact]*Node
	CurrentNode          *Node
}

type Nexus struct {
	Leaders   []*Node           `json:"leaders"`
	cMap      map[uint64]uint64
	Artifacts []*Artifact       `json:"artifacts"`
}

func (n *Nexus) LoadFromFile(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r).Decode(n); err != nil {
		return err
	}

	n.generateCMap()
	return nil
}

func (n *Nexus) generateCMap() {
	// Use a breadth-first search to explore all nodes under each leader.
	n.cMap = make(map[uint64]uint64)
	for _, node := range n.Leaders {
		queue := make([]*Node, 0)
		queue = append(queue, node)
		count := uint64(0)
		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]
			if curr.LeftChild != nil {
				queue = append(queue, curr.LeftChild)
				count++
			}
			if curr.RightChild != nil {
				queue = append(queue, curr.RightChild)
				count++
			}
		}
		n.cMap[node.Id] = count
	}
}

func (n *Nexus) Start(p *Player) {
	// Assign a random leader node to the player.
	p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders))]
	// Initialize Artifact-Distribution
	p.ArtifactDistribution = make(map[*Artifact]*Node)
	n.scrambleArtifacts(p, true)
}

func (n *Nexus) Traverse(p *Player, opt Option) {
	if p.CurrentNode.IsLeaf {
		p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders))]
		n.scrambleArtifacts(p, false)
		return
	}

	p.CurrentNode = p.CurrentNode.Traverse(opt)
}

func (n *Nexus) scrambleArtifacts(p *Player, forceScramble bool) {
	// We scramble the artifacts based on their scramble-coefficients.
	for _, artifact := range n.Artifacts {
		num := rand.Intn(1000)
		if num >= int(1000*artifact.ScrambleCoefficient) || forceScramble {

			leader := p.CurrentNode.FetchLeader()

			toughCoeff := 1 - artifact.ScrambleCoefficient
			nodePath := 1 + int(toughCoeff * float64(rand.Intn(1 + int(n.cMap[leader.Id]))))

			p.ArtifactDistribution[artifact] = p.CurrentNode.GetNodeByNum(nodePath)
		}
	}
}
