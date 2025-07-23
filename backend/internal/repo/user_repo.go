package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kar1mov-u/LeetClone/internal/models"
)

type UserRepository struct {
	conn *pgxpool.Pool
}

func NewUserRepo(conn *pgxpool.Pool) *UserRepository {
	return &UserRepository{conn: conn}
}

func (r *UserRepository) CheckUsername(context context.Context, username string) bool {
	res := r.conn.QueryRow(context, "SELECT email FROM users WHERE username=$1", username)
	email := ""
	err := res.Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true
		}
	}
	return false
}

func (r *UserRepository) CheckEmail(context context.Context, email string) bool {
	res := r.conn.QueryRow(context, "SELECT username FROM users WHERE email=$1", email)
	username := ""
	err := res.Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true
		}
	}
	return false

}

func (r *UserRepository) CreateUser(context context.Context, data models.UserRegister) (uuid.UUID, error) {
	id := uuid.New()
	query := "INSERT INTO users (username, email, password) VALUES($1, $2, $3) RETURNING id"
	rows, err := r.conn.Query(context, query, data.Username, data.Email, data.Password)
	if err != nil {
		return id, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}
