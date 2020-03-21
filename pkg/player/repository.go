package player

import (
	"bandersnatch/pkg/entities"
)

type Repository interface {
	SignUp(p *entities.Player) error
	SignIn(p *entities.Player) error
	Find(email string) (*entities.Player, error)
}
