package game

import (
	"bandersnatch/pkg"
	"encoding/json"
	"math/rand"
	"os"
)

type Player struct {
	Id                   uint64
	CurrentNode          Node
	CollectedArtifacts   map[uint64]struct{}
	TotalScore           uint64
	VisitedNodeMap       map[uint64]struct{}
}

type Nexus struct {
	Leaders       []*Node
	Nodes         []*Node `json:"nodes"`
	Players       map[uint64]*Player
	Artifacts     []*Artifact `json:"artifacts"`
	artifactMap   map[uint64]*Artifact
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
	artifactMap := make(map[uint64]*Artifact)
	for _, artifact := range n.Artifacts {
		artifactMap[artifact.Id] = artifact
	}
	n.artifactMap = artifactMap

	n.artifactNodes = make(map[*Node][]*Node)
	for _, node := range n.Nodes {
		if len(node.ArtifactIds) > 0 {
			if _, ok := n.artifactNodes[node.Leader]; !ok {
				n.artifactNodes[node.Leader] = make([]*Node, 0)
			}
			n.artifactNodes[node.Leader] = append(n.artifactNodes[node.Leader], node)
		}
	}
}

func (n *Nexus) Start(p *Player) error {
	if p == nil {
		return pkg.ErrNilNode
	}
	n.Players[p.Id] = &Player{Id: p.Id}
	*p = *n.Players[p.Id]
	p.VisitedNodeMap = make(map[uint64]struct{})
	// Assign a random leader node to the player.
	p.TotalScore = 0


	p.CurrentNode = *n.Leaders[rand.Intn(len(n.Leaders))]

	p.CollectedArtifacts = make(map[uint64]struct{})
	// Initialize Artifact-Distribution
	*n.Players[p.Id] = *p
	return nil
}

func (n *Nexus) InjectArtifacts(p *Player) {
	// This is done to prevent game-state injection
	*p = *n.Players[p.Id]
	for _, id := range p.CurrentNode.ArtifactIds {
		if _, ok := p.CollectedArtifacts[id]; ok {
			continue
		}
		p.CollectedArtifacts[id] = struct{}{}
		p.TotalScore += n.artifactMap[id].Score
	}
	*n.Players[p.Id] = *p
}

func (n *Nexus) satisfiesDependency(target *Node, p *Player) bool {
	for _, id := range target.RequiredArtifactIds {
		if _, ok := p.CollectedArtifacts[id]; !ok {
			return false
		}
	}
	return true
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
		if len(p.CurrentNode.LeftNodeIds) == 0 && len(p.CurrentNode.RightNodeIds) == 0 {
			p.CurrentNode.IsLeaf = true
		} else {
			if len(p.CurrentNode.LeftNodeIds) == 0 {
				p.CurrentNode.LeftNodeIds = p.CurrentNode.RightNodeIds
			}
			if len(p.CurrentNode.RightNodeIds) == 0 {
				p.CurrentNode.RightNodeIds = p.CurrentNode.LeftNodeIds
			}

			// randomize and prevent subtree selection starvation.
			length := len(p.CurrentNode.LeftNodeIds)
			randId := rand.Intn(length)
			cycleCompleted := false
			endId := (randId+length-1)%length
			for i := randId; ; i = (i + 1) % length {
				p.CurrentNode.LeftChild = n.Nodes[p.CurrentNode.LeftNodeIds[i]-1]

				if !n.satisfiesDependency(p.CurrentNode.LeftChild, p) {
					if i == endId && cycleCompleted {
						// If no viable-node is found, then declare the node as a leaf-node
						p.CurrentNode.IsLeaf = true
						break
						//return pkg.ErrNoPathFound
					}
					if i == endId {
						cycleCompleted = true
					}
					continue
				}

				if i == endId {
					cycleCompleted = true
				}

				if _, ok := p.VisitedNodeMap[p.CurrentNode.LeftChild.Id]; !ok || cycleCompleted {
					break
				}
			}

			length = len(p.CurrentNode.RightNodeIds)
			randId = rand.Intn(length)
			cycleCompleted = false
			endId = (randId+length-1)%length

			for i := randId; ; i = (i + 1) % length {
				p.CurrentNode.RightChild = n.Nodes[p.CurrentNode.RightNodeIds[i]-1]



				if !n.satisfiesDependency(p.CurrentNode.RightChild, p) {
					if i == endId && cycleCompleted {
						// If no viable-node is found, then declare the node as a leaf-node
						p.CurrentNode.IsLeaf = true
						break
						//return pkg.ErrNoPathFound
					}
					if i == endId {
						cycleCompleted = true
					}
					continue
				}

				if i == endId {
					cycleCompleted = true
				}

				if _, ok := p.VisitedNodeMap[p.CurrentNode.RightChild.Id]; !ok || cycleCompleted {
					break
				}
			}
		}
	}

	if p.CurrentNode.IsLeaf {
		p.CurrentNode = *n.Leaders[rand.Intn(len(n.Leaders))]
	} else {
		p.CurrentNode = *p.CurrentNode.Traverse(opt)
	}
	p.VisitedNodeMap[p.CurrentNode.Id] = struct{}{}
	*n.Players[p.Id] = *p
	n.InjectArtifacts(p)

	return nil
}

func (n *Nexus) CheckIfPlayerExists(p *Player) bool {
	if p == nil {
		return false
	}
	_, ok := n.Players[p.Id]
	return ok
}

func (n *Nexus) LoadGameState(p *Player) error {
	if p == nil {
		return pkg.ErrNilNode
	}

	if pl, ok := n.Players[p.Id]; ok {
		*p = *pl
		return nil
	} else {
		return pkg.ErrNotFound
	}
}

func (n *Nexus) FetchArtifactById(id uint64) *Artifact {
	return n.artifactMap[id]
}
