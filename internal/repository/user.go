package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/models"
	"github.com/google/uuid"
)


func (repo *Repository) GetUserByEmail(email string) (models.LoginUser, error) {
	query := `
		SELECT id, email, password_hash, role 
		FROM users WHERE email = $1;
	`
	var user models.LoginUser

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (repo *Repository) RegisterUser(user *models.RegisterUser) (uuid.UUID, error) {
	const query = `
		INSERT INTO users (first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	//TODO Separate each users
	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.PasswordHash == "" {
		return uuid.Nil, errors.New("all fields are required (first Name, last Name, email, password)")
	}

	var id uuid.UUID
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exisitingUser, err := repo.GetUserByEmail(user.Email)
	if err == nil {
		return  uuid.Nil, err
	}

	if exisitingUser.Email != "" {
		return  uuid.Nil, errors.New("user already exists")
	}

	err = repo.db.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (repo *Repository) LoginUser(email string) (models.LoginUser, error) {
	var user models.LoginUser
	var err error

	user, err = repo.GetUserByEmail(email)
	if err != nil {
		return user, err
	}

	return user, nil
}



func (repo *Repository) UserDetails (id uuid.UUID) (models.UserDetails, error) {
	query := `
		SELECT id, first_name, last_name, email, role
		FROM users 
		WHERE id = $1;
	`

	var user models.UserDetails

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role)

	if err != nil {
		return  user, err
	}

	return user, nil
}