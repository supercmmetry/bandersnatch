package core

import (
	"bandersnatch/pkg"
	"math/rand"
)

type GameData struct {
	VisitedNexusPorts  []*Node
	CollectedArtifacts uint
}

type BanderSnatch struct {
	NexusPorts  []*Node
	CurrentNode *Node
	Data        GameData
}

func (b *BanderSnatch) Traverse(option bool) error {
	curr := b.CurrentNode
	state := curr.State
	if state.ContainsArtifact {
		b.Data.CollectedArtifacts++
		curr.State.ContainsArtifact = false
	}

	if curr.IsLeaf {
		if state.LeapToNexus {
			// todo: Implement Bandersnatch Nexus.
			// Skip nexus and warp to a random nexus port.
			b.CurrentNode = b.NexusPorts[rand.Intn(len(b.NexusPorts))]
			b.Data.VisitedNexusPorts = append(b.Data.VisitedNexusPorts, b.CurrentNode)
		} else if state.JumpBack {
			rewindLimit := rand.Intn(BandersnatchRewindLimit)
			for i := 0; i < rewindLimit; i++ {
				if b.CurrentNode.IsLeader {
					break
				}
				b.CurrentNode = b.CurrentNode.Parent
			}
		} else if state.LeapToNode != nil {
			b.CurrentNode = state.LeapToNode
		}
		return nil
	}

	switch option {
	case BandersnatchLeft:
		if b.CurrentNode.LeftChild == nil {
			return pkg.ErrNilChild
		}
		b.CurrentNode = b.CurrentNode.LeftChild
	case BandersnatchRight:
		if b.CurrentNode.RightChild == nil {
			return pkg.ErrNilChild
		}
		b.CurrentNode = b.CurrentNode.RightChild
	}
	return nil
}
