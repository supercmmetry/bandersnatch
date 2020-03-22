package game

import (
	"bandersnatch/pkg"
	"bandersnatch/pkg/entities"
)

type Service struct {
	nexus *Nexus
}

func NewService(n *Nexus) *Service {
	return &Service{nexus: n}
}

func (s *Service) StartGame(p *entities.Player) *NodeData {
	player := &Player{Id: p.Id}
	s.nexus.Start(player)
	return player.CurrentNode.Data
}

func (s *Service) Play(p *Player, opt Option) (*NodeData, error) {
	if opt != OptionLeft && opt != OptionRight {
		return nil, pkg.ErrInvalidOperation
	}
	s.nexus.Traverse(p, opt)
	return p.CurrentNode.Data, nil
}

func (s *Service) GetArtifacts(p *Player) ([]*entities.AbstractArtifact, error) {
	player, ok := s.nexus.Players[p.Id]
	if !ok {
		return nil, pkg.ErrNotFound
	}
	artifacts := make([]*entities.AbstractArtifact, 0)
	for k := range player.CollectedArtifacts {
		artifacts = append(artifacts, &entities.AbstractArtifact{Id: k.Id, Description: k.Description})
	}

	return artifacts, nil
}