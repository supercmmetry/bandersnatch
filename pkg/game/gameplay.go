package game

import (
	"bandersnatch/pkg"
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
	Leaders       []*Node
	Nodes         []*Node `json:"nodes"`
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

	n.createDyraStat()
	n.generateArtifactNodes()
	n.Players = make(map[uint64]*Player)
	return nil
}

func (n *Nexus) createDyraStat() {
	// we iterate through each node in the nexus to generate a dyrastat
	nodeMap := make(map[uint64]*Node)
	for _, node := range n.Nodes {
		nodeMap[node.Id] = node
		if len(node.LeftNodeIds) > 1 || len(node.RightNodeIds) > 1 {
			node.RandomizePath = true
		}
		if len(node.LeftNodeIds) == 0 && len(node.RightNodeIds) == 0 {
			node.IsLeaf = true
		}
	}

	for _, node := range nodeMap {
		if node.IsLeader {
			n.Leaders = append(n.Leaders, node)
			node.Leader = node
		}
		if node.RandomizePath {
			for _, childId := range node.LeftNodeIds {
				child := nodeMap[childId]
				child.Leader = node.Leader
			}
			for _, childId := range node.RightNodeIds {
				child := nodeMap[childId]
				child.Leader = node.Leader
			}
			continue
		}
		if len(node.LeftNodeIds) > 0 && node.LeftNodeIds[0] != 0 {
			node.LeftChild = nodeMap[node.LeftNodeIds[0]]
			node.LeftChild.Leader = node.Leader
		}
		if len(node.RightNodeIds) > 0 && node.RightNodeIds[0] != 0 {
			node.RightChild = nodeMap[node.RightNodeIds[0]]
			node.RightChild.Leader = node.Leader
		}
	}
}

func (n *Nexus) generateArtifactNodes() {
	n.artifactNodes = make(map[*Node][]*Node)
	for _, node := range n.Nodes {
		if node.CanHoldArtifact {
			if _, ok := n.artifactNodes[node.Leader]; !ok {
				n.artifactNodes[node.Leader] = make([]*Node, 0)
			}
			n.artifactNodes[node.Leader] = append(n.artifactNodes[node.Leader], node)
		}
	}
}

func (n *Nexus) Start(p *Player) error {
	n.Players[p.Id] = &Player{Id: p.Id}
	*p = *n.Players[p.Id]
	// Assign a random leader node to the player.
	p.TotalScore = 0
	p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders))]
	p.CollectedArtifacts = make(map[*Artifact]struct{})
	// Initialize Artifact-Distribution
	*n.Players[p.Id] = *p
	return n.scrambleArtifacts(p, true)
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

func (n *Nexus) Traverse(p *Player, opt Option) error {
	if p == nil {
		return pkg.ErrNilNode
	}

	if _, ok := n.Players[p.Id]; !ok {
		return pkg.ErrNilNode
	}

	*p = *n.Players[p.Id]

	if p.CurrentNode.RandomizePath {
		if err := n.scrambleArtifacts(p, false); err != nil {
			return err
		}
		if len(p.CurrentNode.LeftNodeIds) == 0 && len(p.CurrentNode.RightNodeIds) == 0 {
			p.CurrentNode.IsLeaf = true
		} else {
			if len(p.CurrentNode.LeftNodeIds) == 0 {
				p.CurrentNode.LeftNodeIds = p.CurrentNode.RightNodeIds
			}
			if len(p.CurrentNode.RightNodeIds) == 0 {
				p.CurrentNode.RightNodeIds = p.CurrentNode.LeftNodeIds
			}

			p.CurrentNode.LeftChild = n.Nodes[p.CurrentNode.LeftNodeIds[rand.Intn(len(p.CurrentNode.LeftNodeIds))]-1]
			p.CurrentNode.RightChild = n.Nodes[p.CurrentNode.RightNodeIds[rand.Intn(len(p.CurrentNode.RightNodeIds))]-1]
		}
	}

	if p.CurrentNode.IsLeaf {
		p.CurrentNode = n.Leaders[rand.Intn(len(n.Leaders))]
		*n.Players[p.Id] = *p
		if err := n.scrambleArtifacts(p, false); err != nil {
			return err
		}
	} else {
		p.CurrentNode = p.CurrentNode.Traverse(opt)
	}
	*n.Players[p.Id] = *p
	n.CheckForArtifact(p)

	return nil
}

func (n *Nexus) scrambleArtifacts(p *Player, forceScramble bool) error {
	if p == nil {
		return pkg.ErrNilNode
	}

	if _, ok := n.Players[p.Id]; !ok {
		return pkg.ErrNilNode
	}

	*p = *n.Players[p.Id]
	p.ArtifactDistribution = make(map[*Node]*Artifact)
	// We scramble the artifacts based on their scramble-coefficients.
	for _, artifact := range n.Artifacts {
		if _, ok := p.CollectedArtifacts[artifact]; ok {
			continue
		}
		num := rand.Intn(1000)
		if num >= int(1000*artifact.ScrambleCoefficient) || forceScramble {
			leader := p.CurrentNode.Leader
			anodes := n.artifactNodes[leader]
			anode := anodes[rand.Intn(len(anodes))]
			p.ArtifactDistribution[anode] = artifact
			break
		}
	}
	*n.Players[p.Id] = *p
	return nil
}
