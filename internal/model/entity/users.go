package entity

import "time"

type Users struct {
	UserID    int64     `gorm:"primaryKey" json:"user_id"`
	Username  string    `gorm:"type:varchar(55)" json:"username"`
	Password  string    `json:"password"`
	Email     string    `gorm:"unique" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
