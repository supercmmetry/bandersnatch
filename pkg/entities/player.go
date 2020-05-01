package entities

type Player struct {
	Id       uint64 `json:"id" gorm:"auto_increment"`
	Email    string `json:"email" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
	MaxScore uint64 `json:"score"`
}

type AbstractArtifact struct {
	Id            uint64                 `json:"id"`
	Miscellaneous map[string]interface{} `json:"misc"`
}

type AbstractPlayer struct {
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"primary_key"`
	MaxScore uint64 `json:"score"`
}
