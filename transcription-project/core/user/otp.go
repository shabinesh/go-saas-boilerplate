package user

import "time"

type OTP struct {
	ID        int
	UserID    string
	Code      string
	CreatedAt time.Time
}
