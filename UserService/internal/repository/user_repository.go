package repository

import (
	"database/sql"
	"log"
)

type UserRepository interface {
	CreateUser(name, surname, phoneNumber, role string) (int32, string, error)
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(name, surname, phoneNumber, role string) (int32, string, error) {
	_, err := r.db.Exec("INSERT INTO users (name, surname, phone_number, role) VALUES ($1, $2, $3, $4)", name, surname, phoneNumber, role)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 500, "Internal server error", err
	}
	return 200, "User created successfully", nil
}
