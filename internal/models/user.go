package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"size:255"`
	Email    string `json:"email" gorm:"unique;size:255"`
	Password string `json:"-" gorm:"size:255"`
}
