package repo

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/shabinesh/transcription/core/user"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) FindUser(id string) (*user.User, bool, error) {
	r := u.db.QueryRow("SELECT * FROM users WHERE id = ?", id)

	var uu user.User
	err := r.Scan(&uu.ID, &uu.Email, &uu.IsVerified, &uu.IsActive, &uu.CreatedAt)
	if err != nil {
		return nil, false, err
	}

	return &uu, true, nil
}

func (u *UserRepo) FindUserByEmail(email string) (*user.User, bool, error) {
	r := u.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	var uu user.User
	err := r.Scan(uu.ID, uu.Email, uu.IsVerified, uu.IsActive, uu.CreatedAt)
	if err != nil {
		return nil, false, err
	}

	return &uu, true, nil
}

func (u *UserRepo) AddUser(uu *user.User) (*user.User, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec("INSERT INTO users (id, email, is_verified, is_active) VALUES (?, ?, ?, ?)", uid.String(), uu.Email, uu.IsVerified, uu.IsActive)
	if err != nil {
		return nil, err
	}

	uu.ID = user.UserID(uid.String())

	return uu, nil
}

func (u *UserRepo) UpdateUserStatus(uu *user.User) error {
	_, err := u.db.Exec("UPDATE users SET is_verified = ?, is_active = ? WHERE id = ?", uu.IsVerified, uu.IsActive, uu.ID)
	if err != nil {
		return err
	}

	return nil
}
