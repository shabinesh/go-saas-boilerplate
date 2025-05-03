package otp

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"time"

	"github.com/shabinesh/app/core/user"
)

const maxDigits = 6

type OTPRepository interface {
	SaveOTP(userID user.UserID, otp string, reason string) error
	GetOTP(userID user.UserID, reason string) (*user.OTP, error)
}

type OTPService struct {
	OTPRepo OTPRepository
}

func NewOTPService(repo OTPRepository) *OTPService {
	return &OTPService{OTPRepo: repo}
}

func (a *OTPService) Generate(userID string, reason string) string {
	k, _ := rand.Int(rand.Reader, big.NewInt(int64(math.Pow(10, float64(maxDigits)))))
	otp := fmt.Sprintf("%0*d", maxDigits, k)
	err := a.OTPRepo.SaveOTP(user.UserID(userID), otp, reason)
	if err != nil {
		slog.Error("failed to save otp", err)
		return ""
	}

	return otp
}

func (a *OTPService) Verify(userID string, code string, reason string) (bool, error) {
	otp, err := a.OTPRepo.GetOTP(user.UserID(userID), reason)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("No otp found", "userID", userID, "reason", reason)
			return false, fmt.Errorf("invalid otp")
		}
		return false, err
	}

	if otp.Code == "" {
		slog.Error("No otp found", "userID", userID, "reason", reason)
		return false, fmt.Errorf("invalid otp")
	}

	if otp.CreatedAt.Before(time.Now().Add(-10 * time.Minute)) {
		slog.Error("OTP expired", "userID", userID, "reason", reason)
		return false, fmt.Errorf("otp expired")
	}

	if otp.Code != code {
		slog.Error("Invalid doesn't match", "userID", userID, "reason", reason, "code", code, "otp", otp.Code)
		return false, nil
	}

	return true, nil
}
