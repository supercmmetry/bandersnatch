package game

import (
	"bandersnatch/pkg"
	"encoding/json"
	"math/rand"
	"os"
)

type Player struct {
	Id                 uint64
	CurrentNode        Node
	CollectedArtifacts map[uint64]struct{}
	TotalScore         uint64
	CycleMap           map[uint64]struct{}
	VisitedNodeMap     map[uint64]struct{}
	HintIdx            uint8
}

type HintKey struct {
	NodeId  uint64
	HintIdx uint8
}

type Nexus struct {
	Leaders       []*Node
	Nodes         []*Node `json:"nodes"`
	Players       map[uint64]*Player
	Artifacts     []*Artifact `json:"artifacts"`
	artifactMap   map[uint64]*Artifact
	artifactNodes map[*Node][]*Node // maps a leader node to a list of potential artifact nodes under the leader node.
	HintMap       map[HintKey]Hint
}

type HintData struct {
	Hints []*Hint `json:"hints"`
}

func (h *HintData) LoadFromFile(filename string) error {
	r, err := os.Open(filename)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(r).Decode(h); err != nil {
		return err
	}

	return nil
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

func (n *Nexus) LoadHintsFromFile(filename string) error {
	hintData := HintData{}
	if err := hintData.LoadFromFile(filename); err != nil {
		return err
	}

	n.HintMap = make(map[HintKey]Hint)

	for _, hint := range hintData.Hints {
		n.HintMap[HintKey{hint.NodeID, hint.HintIdx}] = *hint

	}

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
	p.CycleMap = make(map[uint64]struct{})
	p.VisitedNodeMap = make(map[uint64]struct{})
	// Assign a random leader node to the player.
	p.TotalScore = 0

	p.CurrentNode = *n.Leaders[rand.Intn(len(n.Leaders))]
	p.CycleMap[p.CurrentNode.Id] = struct{}{}
	p.VisitedNodeMap[p.CurrentNode.Id] = struct{}{}

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

	// Check for anti-requisite artifacts
	for _, id := range target.AntiRequisiteArtifactIds {
		if _, ok := p.CollectedArtifacts[id]; ok {
			return false
		}
	}

	return true
}

func (n *Nexus) cyclicRefresh(p *Player) {
	leftCycleCompleted := true
	rightCycleCompleted := true

	nodeIds := p.CurrentNode.LeftNodeIds
	for _, id := range nodeIds {
		if _, ok := p.CycleMap[id]; !ok && n.satisfiesDependency(n.Nodes[id-1], p) {
			leftCycleCompleted = false
			break
		}
	}

	nodeIds = p.CurrentNode.RightNodeIds
	for _, id := range nodeIds {
		if _, ok := p.CycleMap[id]; !ok && n.satisfiesDependency(n.Nodes[id-1], p) {
			rightCycleCompleted = false
			break
		}
	}

	// This means that we have visited all nodes in LRU, and we need to refresh them to counter
	// starvation in the next pass, otherwise we would go to a random-node in every pass which could
	// create starvation.
	// For this we "unvisit" the candidate nodes.

	if leftCycleCompleted {
		for _, id := range p.CurrentNode.LeftNodeIds {
			delete(p.CycleMap, id)
		}
	}
	if rightCycleCompleted {
		for _, id := range p.CurrentNode.RightNodeIds {
			delete(p.CycleMap, id)
		}
	}

	// Now we might have removed the currentNode from the LRU.
	// So we resolve it by re-adding the current node in the LRU

	p.CycleMap[p.CurrentNode.Id] = struct{}{}
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
		// refresh LRU-cache to minimize starvation for next-pass.
		n.cyclicRefresh(p)

		if len(p.CurrentNode.LeftNodeIds) == 0 && len(p.CurrentNode.RightNodeIds) == 0 {
			p.CurrentNode.IsLeaf = true
		} else {
			// left-right node resolution
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
			endId := (randId + length - 1) % length
			for i := randId; ; i = (i + 1) % length {
				p.CurrentNode.LeftChild = n.Nodes[p.CurrentNode.LeftNodeIds[i]-1]

				if !n.satisfiesDependency(p.CurrentNode.LeftChild, p) {
					if i == endId && cycleCompleted {
						// If no viable-node is found, then declare the node as a leaf-node
						p.CurrentNode.IsLeaf = true
						break
					}
					if i == endId {
						cycleCompleted = true
					}
					continue
				}

				if i == endId {
					cycleCompleted = true
				}

				if _, ok := p.CycleMap[p.CurrentNode.LeftChild.Id]; !ok || cycleCompleted {
					break
				}
			}

			length = len(p.CurrentNode.RightNodeIds)
			randId = rand.Intn(length)
			cycleCompleted = false
			endId = (randId + length - 1) % length

			for i := randId; ; i = (i + 1) % length {
				p.CurrentNode.RightChild = n.Nodes[p.CurrentNode.RightNodeIds[i]-1]

				if !n.satisfiesDependency(p.CurrentNode.RightChild, p) {
					if i == endId && cycleCompleted {
						// If no viable-node is found, then declare the node as a leaf-node
						p.CurrentNode.IsLeaf = true
						break
					}
					if i == endId {
						cycleCompleted = true
					}
					continue
				}

				if i == endId {
					cycleCompleted = true
				}

				if _, ok := p.CycleMap[p.CurrentNode.RightChild.Id]; !ok || cycleCompleted {
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

	// Increment total score by one if a new node is visited.
	// Reset HintIdx to 0
	if _, ok := p.VisitedNodeMap[p.CurrentNode.Id]; !ok {
		p.TotalScore += 1
	}

	p.HintIdx = 0

	p.CycleMap[p.CurrentNode.Id] = struct{}{}
	p.VisitedNodeMap[p.CurrentNode.Id] = struct{}{}

	*n.Players[p.Id] = *p
	n.InjectArtifacts(p)

	return nil
}

func (n *Nexus) GetHint(p *Player) (*Hint, error) {
	player := n.Players[p.Id]

	hint, ok := n.HintMap[HintKey{player.CurrentNode.Id, player.HintIdx}]
	if !ok {
		return nil, pkg.ErrNoHintFound
	}

	if player.TotalScore < hint.Penalty {
		return nil, pkg.ErrInsufficientScore
	} else {
		player.TotalScore -= hint.Penalty
		player.HintIdx += 1
	}

	*n.Players[p.Id] = *player

	return &hint, nil
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
