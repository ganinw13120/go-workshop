package entity

import "time"

type User struct {
	ID        string    `example:"1234"`
	Name      string    `example:"John Doe"`
	Username  string    `example:"johndoe"`
	Password  string    `example:"A1b2C3d$"`
	BirthDate time.Time `json:"birth_date" example:"2006-01-02"`
}

type UserSession struct {
	UserID   string `mapstructure:"user_id" json:"user_id"`
	Username string `mapstructure:"username" json:"username"`
}
