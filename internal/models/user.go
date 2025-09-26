package models

import (
	"context"
	"errors"
	"fmt"
	"jwt-golang/internal/common"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id             string `db:"id"`
	Username       string `db:"username"`
	HashedPassword string `db:"hashed_password"`
	Email          string `db:"email"`
}

type PostUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserPassword struct {
	Id             string `db:"id"`
	HashedPassword string `db:"hashed_password"`
}

type UserLogin struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type EditUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	NewPassword string `json:"new_password,omitempty"` // optional
}

type UserModelInterface interface {
	Insert(ctx context.Context, userInput PostUser) (string, error)
	GetAll(ctx context.Context) ([]Users, error)
	AuthenticateUser(ctx context.Context, userLoginData UserLogin) (*Users, error)
	GetUserById(ctx context.Context, userId string) (*Users, error)
	EditUser(ctx context.Context, userInput EditUser, currentEmail string) error
	DeleteUser(ctx context.Context, email string) error
}

type UserModel struct {
	Pool *pgxpool.Pool
}

func (m *UserModel) Insert(ctx context.Context, userInput PostUser) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := m.Pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	id := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), 12)
	if err != nil {
		return "", err
	}

	query := `INSERT INTO users (username, hashed_password, id, email) VALUES(@username, @hashed_password, @id ,@email)`
	args := pgx.NamedArgs{
		"username":        userInput.Username,
		"hashed_password": hashedPassword,
		"id":              id,
		"email":           userInput.Email,
	}

	commandTag, err := tx.Exec(ctx, query, args)
	if err != nil {
		return "", err
	}

	if commandTag.RowsAffected() != 1 {
		return "", errors.New("failed to insert user")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (m *UserModel) GetAll(ctx context.Context) ([]Users, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `SELECT * FROM users`

	rows, err := m.Pool.Query(ctx, query)
	if err != nil {
		return nil, common.ErrQueryError
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[Users])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) GetUserById(ctx context.Context, userId string) (*Users, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `SELECT * FROM users where id = @id`
	args := pgx.NamedArgs{
		"id": userId,
	}

	rows, err := m.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Users])
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) AuthenticateUser(ctx context.Context, userLoginData UserLogin) (*Users, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := `SELECT * FROM users where email = @email`
	args := pgx.NamedArgs{
		"email": userLoginData.Email,
	}

	rows, err := m.Pool.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Users])
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userLoginData.Password))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) DeleteUser(ctx context.Context, email string) error {
	query := `DELETE FROM users where email = @email`
	args := pgx.NamedArgs{
		"email": email,
	}

	commandTag, err := m.Pool.Exec(ctx, query, args)
	if commandTag.RowsAffected() != 1 {
		return err
	}

	if err != nil {
		return err
	}

	return err
}

func (m *UserModel) EditUser(ctx context.Context, userInput EditUser, currentUserId string) error {
	// build dynamic query based on what's being updated
	setParts := []string{}
	args := pgx.NamedArgs{"current_id": currentUserId}

	if userInput.Username != "" {
		setParts = append(setParts, "username = @username")
		args["username"] = userInput.Username
	}

	if userInput.Email != "" {
		setParts = append(setParts, "email = @email")
		args["email"] = userInput.Email
	}

	if userInput.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInput.NewPassword), 12)
		if err != nil {
			return err
		}
		setParts = append(setParts, "hashed_password = @hashed_password")
		args["hashed_password"] = hashedPassword
	}

	if len(setParts) == 0 {
		return errors.New("no fields to update")
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE email = @current_email", strings.Join(setParts, ", "))

	commandTag, err := m.Pool.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return errors.New("user not found or no changes made")
	}

	return nil
}
