package player

import (
	"bandersnatch/api/middleware"
	"bandersnatch/pkg/entities"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repo
}

func (s *Service) SignUp(p *entities.Player) (*jwt.Token, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	p.Password = string(passHash)
	token := middleware.Token{Email: p.Email, Password: p.Password}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, token)

	err = s.repo.SignUp(p)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (s *Service) SignIn(p *entities.Player) (*jwt.Token, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	p.Password = string(passHash)
	token := middleware.Token{Email: p.Email, Password: p.Password}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, token)

	return tk, s.repo.SignIn(p)
}

func (s *Service) Find(email string) (*entities.Player, error) {
	return s.repo.Find(email)
}

func NewService(r repo) Service {
	return Service{repo: r}
}
