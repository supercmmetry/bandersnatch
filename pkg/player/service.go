package player

import "bandersnatch/pkg/entities"

type Service struct {
	repo repo
}

func (s *Service) AddPlayer(p *entities.Player) error {
	return s.repo.AddPlayer(p)
}

func (s *Service) Find(email string) (*entities.Player, error) {
	return s.repo.Find(email)
}

func NewService(r repo) Service {
	return Service{repo: r}
}

