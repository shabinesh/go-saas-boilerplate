package repo

import (
	"context"
	"database/sql"
	"errors"
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

func (o *OTPRepository) SaveOTP(userID user.UserID, code string, reason string) error {
	_, err := o.db.Exec(context.Background(), `INSERT INTO otps (user_id, otp_code, reason	)
	VALUES ($1, $2, $3)
	ON CONFLICT(user_id) DO UPDATE SET otp_code = $2, reason = $3;`, userID, code, reason)
	if err != nil {
		return err
	}

	return nil
}

func (o *OTPRepository) GetOTP(userID user.UserID, reason string) (*user.OTP, error) {
	var otp user.OTP
	err := o.db.QueryRow(context.Background(), "SELECT id, user_id, otp_code, created_at FROM otps WHERE user_id = $1 AND reason = $2", userID, reason).Scan(&otp.ID, &otp.UserID, &otp.Code, &otp.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid otp")
		}

		return nil, err
	}

	return &otp, nil
}
