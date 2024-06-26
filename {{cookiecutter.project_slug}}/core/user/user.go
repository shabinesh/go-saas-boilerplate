package user

import "time"

type UserID string

type User struct {
	ID         UserID
	Email      string
	IsVerified bool
	IsActive   bool
	CreatedAt  time.Time
	Meta       map[string]string
}
