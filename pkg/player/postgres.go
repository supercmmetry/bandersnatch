package player

import (
	"bandersnatch/pkg"
	"bandersnatch/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func (r *repo) AddPlayer(player *entities.Player) error {
	tx := r.DB.Begin()

	if err := tx.Where("email = ?", player.Email).Find(&entities.Player{}).Error; err == nil {
		return pkg.ErrAlreadyExists
	} else if err == gorm.ErrRecordNotFound {
		if err := tx.Save(player).Error; err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	} else {
		return err
	}

	return nil
}

func (r *repo) Find(email string) (*entities.Player, error) {
	tx := r.DB.Begin()
	p := &entities.Player{}

	if err := tx.Where("email = ?", email).Find(p).Error; err != nil {
		return nil, err
	}

	return p, nil
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}