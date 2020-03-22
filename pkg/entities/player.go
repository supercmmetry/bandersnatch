package entities

type Player struct {
	Id       uint64 `json:"id" gorm:"auto_increment"`
	Email    string `json:"email" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

type AbstractArtifact struct {
	Id                  uint64  `json:"id"`
	Description         string  `json:"description"`
}
