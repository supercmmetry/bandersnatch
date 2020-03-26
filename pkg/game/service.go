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

func (s *Service) StartGame(p *entities.Player) (*NodeData, error) {
	player := &Player{Id: p.Id}
	if err := s.nexus.Start(player); err != nil {
		return nil, err
	}
	return player.CurrentNode.Data, nil
}

func (s *Service) Play(p *Player, opt Option) (*NodeData, error) {
	if opt != OptionLeft && opt != OptionRight {
		return nil, pkg.ErrInvalidOperation
	}
	if err := s.nexus.Traverse(p, opt); err != nil {
		return nil, err
	}
	return p.CurrentNode.Data, nil
}

func (s *Service) GetNodeData(p *Player) (*NodeData, error) {
	if err := s.nexus.LoadGameState(p); err != nil {
		return nil, err
	}
	return p.CurrentNode.Data, nil
}

func (s *Service) GetArtifacts(p *Player) ([]*entities.AbstractArtifact, error) {
	player, ok := s.nexus.Players[p.Id]
	if !ok {
		return nil, pkg.ErrNotFound
	}
	artifacts := make([]*entities.AbstractArtifact, 0)
	for id := range player.CollectedArtifacts {
		k := s.nexus.FetchArtifactById(id)
		artifacts = append(artifacts, &entities.AbstractArtifact{Id: k.Id, Description: k.Description})
	}

	return artifacts, nil
}

func (s *Service) CheckIfPlayerExists(p *Player) bool {
	return s.nexus.CheckIfPlayerExists(p)
}
