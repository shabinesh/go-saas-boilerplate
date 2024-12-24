package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"

	"github.com/shabinesh/app/core/user"
)

type UserRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) FindUser(id string) (*user.User, bool, error) {
	r := u.db.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1", id)

	var uu user.User
	err := r.Scan(&uu.ID, &uu.Email, &uu.IsVerified, &uu.IsActive, &uu.CreatedAt)
	if err != nil {
		return nil, false, err
	}

	return &uu, true, nil
}

func (u *UserRepo) FindUserByEmail(email string) (*user.User, bool, error) {
	r := u.db.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1", email)
	var uu user.User
	err := r.Scan(&uu.ID, &uu.Email, &uu.IsVerified, &uu.IsActive, &uu.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}

		return nil, false, err
	}

	return &uu, true, nil
}

func (u *UserRepo) AddUser(uu *user.User) (*user.User, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec(context.Background(), "INSERT INTO users (id, email, is_verified, is_active) VALUES ($1, $2, $3, $4)", uid.String(), uu.Email, uu.IsVerified, uu.IsActive)
	if err != nil {
		return nil, err
	}

	uu.ID = user.UserID(uid.String())

	return uu, nil
}

func (u *UserRepo) UpdateUserStatus(uu *user.User) error {
	_, err := u.db.Exec(context.Background(), "UPDATE users SET is_verified = $1, is_active = $2  WHERE id = $3", uu.IsVerified, uu.IsActive, uu.ID)
	if err != nil {
		return err
	}

	return nil
}
