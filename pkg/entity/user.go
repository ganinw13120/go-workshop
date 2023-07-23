package entity

type UserSession struct {
	UserID   string `mapstructure:"user_id" json:"user_id"`
	Username string `mapstructure:"username" json:"username"`
}
