package game

import "bandersnatch/pkg/entities"

type Service struct {
	nexus Nexus
}

func (s *Service) StartGame(p *entities.Player) *Player {
	player := &Player{Id: p.Id}
	s.nexus.Start(player)
	return player
}

func (s *Service) Play(p *Player, opt Option) *NodeData {
	s.nexus.Traverse(p, opt)
	return p.CurrentNode.Data
}
