package repository

import (
	"database/sql"
	"log"
)

type BookRepository interface {
	CreateBook(name string) (int32, string, error)
}

type PostgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{db: db}
}

func (r *PostgresBookRepository) CreateBook(name string) (int32, string, error) {
	_, err := r.db.Exec("INSERT INTO book (title) VALUES ($1)", name)
	if err != nil {
		log.Printf("Error creating book: %v", err)
		return 500, "Internal server error", err
	}
	return 200, "Book created successfully", nil
}
