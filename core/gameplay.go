package bandersnatch

import "math/rand"

type Player struct {
	Id                   uint64
	ArtifactDistribution map[*Artifact]*Node
	CurrentNode          *Node
}

type Nexus struct {
	Leaders   []*Node
	CMap      map[uint64]uint64
	Artifacts []*Artifact
}

func (n *Nexus) Start(p *Player) {
	// Assign a random leader node to the player.
	p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders)-1)]
	// Initialize Artifact-Distribution
	n.scrambleArtifacts(p)
}

func (n *Nexus) Traverse(p *Player, opt Option) {
	p.CurrentNode.Traverse(opt)
}

func (n *Nexus) scrambleArtifacts(p *Player) {
	// We scramble the artifacts based on their scramble-coefficients.
	for _, artifact := range n.Artifacts {
		num := rand.Intn(1000)
		if num >= int(1000 * artifact.ScrambleCoefficient) {

			leader := p.CurrentNode.FetchLeader()

			toughCoeff := 1 - artifact.ScrambleCoefficient
			nodePath := int(toughCoeff * float64(rand.Intn(int(n.CMap[leader.Id]))))

			p.ArtifactDistribution[artifact] = p.CurrentNode.GetNodeByNum(nodePath)
		}
	}
}
