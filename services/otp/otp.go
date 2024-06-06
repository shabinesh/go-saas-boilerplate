package otp

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"math/big"

	"github.com/shabinesh/transcription/core/user"
)

const maxDigits = 6

type OTPRepository interface {
	SaveOTP(userID user.UserID, otp string) error
	GetOTP(userID user.UserID) (*user.OTP, error)
}

type OTPService struct {
	OTPRepo OTPRepository
}

func NewOTPService(repo OTPRepository) *OTPService {
	return &OTPService{OTPRepo: repo}
}

func (a *OTPService) Generate(userID string) string {
	k, _ := rand.Int(rand.Reader, big.NewInt(int64(math.Pow(10, float64(maxDigits)))))
	otp := fmt.Sprintf("%0*d", maxDigits, k)
	err := a.OTPRepo.SaveOTP(user.UserID(userID), otp)
	if err != nil {
		slog.Error("failed to save otp", err)
		return ""
	}

	return otp
}

func (a *OTPService) Verify(userID string, code string) (bool, error) {
	otp, err := a.OTPRepo.GetOTP(user.UserID(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("invalid otp")
		}
		return false, err
	}

	// TODO - Add expiry check
	fmt.Println("compare", otp.Code, code)
	if otp.Code != code {
		return false, nil
	}

	return true, nil
}
