package player

import (
	"bandersnatch/pkg"
	"bandersnatch/pkg/entities"
	"github.com/jinzhu/gorm"
)

type repo struct {
	DB *gorm.DB
}

func (r *repo) SignUp(player *entities.Player) error {
	tx := r.DB.Begin()

	if err := tx.Where("email = ?", player.Email).Find(player).Error; err == nil {
		tx.Rollback()
		return pkg.ErrAlreadyExists
	} else if err == gorm.ErrRecordNotFound {
		if err := tx.Save(player).Error; err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
	} else {
		tx.Rollback()
		return err
	}

	p, err := r.Find(player.Email)
	if err != nil {
		return err
	}
	*player = *p

	return nil
}

func (r *repo) SaveMaxScore(p *entities.Player) error {
	tx := r.DB.Begin()
	prevP := &entities.Player{}
	if err := tx.Where("email = ?", p.Email).Find(prevP).Error; err != nil {
		tx.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			return pkg.ErrNotFound
		default:
			return pkg.ErrDatabase
		}
	}

	if p.MaxScore > prevP.MaxScore {
		if err := tx.Save(p).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (r *repo) Find(email string) (*entities.Player, error) {
	tx := r.DB.Begin()
	p := &entities.Player{}

	if err := tx.Where("email = ?", email).Find(p).Error; err != nil {
		tx.Rollback()
		return nil, pkg.ErrNotFound
	}

	tx.Commit()
	return p, nil
}

func (r *repo) ViewLeaderboard() ([]entities.AbstractPlayer, error) {
	tx := r.DB.Begin()
	var scores []entities.AbstractPlayer

	if rows, err := tx.Model(&entities.Player{}).Select("name, email, max_score").Order("max_score desc").Rows(); err == nil {
		for rows.Next() {
			var ap entities.AbstractPlayer
			if err := rows.Scan(&ap.Name, &ap.Email, &ap.MaxScore); err != nil {
				return nil, pkg.ErrDatabase
			}
			scores = append(scores, ap)
		}
	} else {
		return nil, pkg.ErrDatabase
	}
	
	return scores, nil
}

func NewPostgresRepo(db *gorm.DB) Repository {
	return &repo{DB: db}
}
