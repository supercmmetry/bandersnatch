package player

import (
	"bandersnatch/api/middleware"
	"bandersnatch/pkg"
	"bandersnatch/pkg/entities"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func (s *Service) SaveMaxScore(p *entities.Player) error {
	return s.repo.SaveMaxScore(p)
}

func (s *Service) SignUp(p *entities.Player) (*jwt.Token, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	p.Password = string(passHash)
	err = s.repo.SignUp(p)
	if err != nil {
		return nil, err
	}

	token := middleware.Token{Email: p.Email, Id: p.Id}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, token)

	return tk, nil
}

func (s *Service) SignIn(p *entities.Player) (*jwt.Token, error) {
	player, err := s.repo.Find(p.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(p.Password)); err != nil {
		return nil, pkg.ErrNotFound
	}

	token := middleware.Token{Email: player.Email, Id: player.Id}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS512, token)

	return tk, nil
}

func (s *Service) Find(email string) (*entities.Player, error) {
	return s.repo.Find(email)
}

func (s *Service) ViewLeaderboard() ([]entities.AbstractPlayer, error) {
	return s.repo.ViewLeaderboard()
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}
