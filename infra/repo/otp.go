package repo

import (
	"database/sql"

	"github.com/shabinesh/transcription/core/user"
)

type otpRepository struct {
	db *sql.DB
}

func NewOTPRepository(db *sql.DB) *otpRepository {
	return &otpRepository{db: db}
}

func (o *otpRepository) SaveOTP(userID user.UserID, otp string) error {
	_, err := o.db.Exec(`INSERT INTO otps (user_id, otp_code)
	VALUES (?, ?)
	ON CONFLICT(user_id) DO UPDATE SET otp_code = ?;`, userID, otp, otp)
	if err != nil {
		return err
	}

	return nil
}

func (o *otpRepository) GetOTP(userID user.UserID) (*user.OTP, error) {
	var otp user.OTP
	err := o.db.QueryRow("SELECT id, user_id, otp_code, created_at FROM otps WHERE user_id = ?", userID).Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &otp, nil
}
