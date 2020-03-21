package entity

type Player struct {
	Id uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Email string `json:"email" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`
}
