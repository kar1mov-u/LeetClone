// Repo
package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/kar1mov-u/LeetClone/internal/models"
)

type UserRepository struct {
	conn Queryer
}

func NewUserRepo(q Queryer) *UserRepository {
	return &UserRepository{conn: q}
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
	query := "INSERT INTO users (username, email, password, role) VALUES($1, $2, $3, 'user') RETURNING id"
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

func (r *UserRepository) GetUserPasswordAndRole(context context.Context, username string) (uuid.UUID, string, string, error) {
	var dbPass string
	var role string
	var userID uuid.UUID
	query := `SELECT id,password, role FROM users WHERE username=$1 OR email=$1`
	row := r.conn.QueryRow(context, query, username)
	err := row.Scan(&userID, &dbPass, &role)
	if err != nil {
		return uuid.UUID{}, "", "", err
	}
	return userID, dbPass, role, nil
}

func (r *UserRepository) GetUserByID(context context.Context, id uuid.UUID) (models.User, error) {
	dbUser := models.User{}
	query := `SELECT username, email,role, created_at FROM users where id=$1`
	row := r.conn.QueryRow(context, query, id)
	err := row.Scan(
		&dbUser.Username,
		&dbUser.Email,
		&dbUser.Role,
		&dbUser.Created_at,
	)
	if err != nil {
		log.Println(err)

		return dbUser, err
	}
	return dbUser, nil
}
