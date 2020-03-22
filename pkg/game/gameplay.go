package game

import (
	"encoding/json"
	"math/rand"
	"os"
)

type Player struct {
	Id                   uint64
	ArtifactDistribution map[*Node]*Artifact
	CurrentNode          *Node
	CollectedArtifacts   map[*Artifact]struct{}
	TotalScore           uint64
}

type Nexus struct {
	Leaders       []*Node `json:"leaders"`
	Players       map[uint64]*Player
	Artifacts     []*Artifact       `json:"artifacts"`
	artifactNodes map[*Node][]*Node // maps a leader node to a list of potential artifact nodes under the leader node.
}

func (n *Nexus) LoadFromFile(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r).Decode(n); err != nil {
		return err
	}

	n.generateArtifactNodes()
	n.Players = make(map[uint64]*Player)
	return nil
}

func (n *Nexus) generateArtifactNodes() {
	n.artifactNodes = make(map[*Node][]*Node)
	for _, node := range n.Leaders {
		queue := make([]*Node, 0)
		queue = append(queue, node)
		for len(queue) > 0 {
			curr := queue[0]

			if curr.CanHoldArtifact {
				n.artifactNodes[node] = append(n.artifactNodes[node], curr)
			}

			queue = queue[1:]
			if curr.LeftChild != nil {
				queue = append(queue, curr.LeftChild)
			}
			if curr.RightChild != nil {
				queue = append(queue, curr.RightChild)
			}
		}
	}
}

func (n *Nexus) Start(p *Player) {
	n.Players[p.Id] = &Player{Id: p.Id}
	*p = *n.Players[p.Id]
	// Assign a random leader node to the player.
	p.TotalScore = 0
	p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders))]
	p.CollectedArtifacts = make(map[*Artifact]struct{})
	// Initialize Artifact-Distribution
	*n.Players[p.Id] = *p
	n.scrambleArtifacts(p, true)
}

func (n *Nexus) CheckForArtifact(p *Player) *Artifact {
	// This is done to prevent game-state injection
	*p = *n.Players[p.Id]
	if artifact, ok := p.ArtifactDistribution[p.CurrentNode]; ok {
		p.CollectedArtifacts[artifact] = struct{}{}
		p.TotalScore += uint64(100 * artifact.ScrambleCoefficient)
		*n.Players[p.Id] = *p
		return artifact
	} else {
		return nil
	}
}

func (n *Nexus) Traverse(p *Player, opt Option) {
	*p = *n.Players[p.Id]
	if p.CurrentNode.IsLeaf {
		p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders))]
		*n.Players[p.Id] = *p
		n.scrambleArtifacts(p, false)
	} else {
		p.CurrentNode = p.CurrentNode.Traverse(opt)
	}
	*n.Players[p.Id] = *p
	n.CheckForArtifact(p)
}

func (n *Nexus) scrambleArtifacts(p *Player, forceScramble bool) {
	*p = *n.Players[p.Id]
	p.ArtifactDistribution = make(map[*Node]*Artifact)
	// We scramble the artifacts based on their scramble-coefficients.
	for _, artifact := range n.Artifacts {
		if _, ok := p.CollectedArtifacts[artifact]; ok {
			continue
		}
		num := rand.Intn(1000)
		if num >= int(1000*artifact.ScrambleCoefficient) || forceScramble {
			leader := p.CurrentNode.FetchLeader()
			anodes := n.artifactNodes[leader]
			anode := anodes[rand.Intn(len(anodes))]
			p.ArtifactDistribution[anode] = artifact
			break
		}
	}
	*n.Players[p.Id] = *p
}
