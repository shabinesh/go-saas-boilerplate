package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"

	"github.com/shabinesh/app/core/user"
)

type OTPRepository struct {
	db *pgx.Conn
}

func NewOTPRepository(db *pgx.Conn) *OTPRepository {
	return &OTPRepository{db: db}
}

func (o *OTPRepository) SaveOTP(userID user.UserID, otp string) error {
	_, err := o.db.Exec(context.Background(), `INSERT INTO otps (user_id, otp_code)
	VALUES ($1, $2)
	ON CONFLICT(user_id) DO UPDATE SET otp_code = $3;`, userID, otp, otp)
	if err != nil {
		return err
	}

	return nil
}

func (o *OTPRepository) GetOTP(userID user.UserID) (*user.OTP, error) {
	var otp user.OTP
	err := o.db.QueryRow(context.Background(), "SELECT id, user_id, otp_code, created_at FROM otps WHERE user_id = ?", userID).Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid otp")
		}

		return nil, err
	}

	return &otp, nil
}
