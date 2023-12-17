package domain

import "time"

type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"create_at"`
}
