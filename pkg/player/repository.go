package player

import (
	"bandersnatch/pkg/entities"
)

type Repository interface {
	AddPlayer(p *entities.Player) error
	Find(email string) (*entities.Player, error)
}
