package player

import (
	"bandersnatch/pkg/entities"
)

type Repository interface {
	SignUp(p *entities.Player) error
	Find(email string) (*entities.Player, error)
	SaveMaxScore(p *entities.Player) error
	ViewLeaderboard() ([]entities.AbstractPlayer, error)
}
