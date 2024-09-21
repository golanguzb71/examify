package repository

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"ielts-service/internal/models"
	"log"
	"strconv"
)

type PostgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{db: db}
}

func (r *PostgresBookRepository) CreateBook(name string) error {
	var checker bool
	err := r.db.QueryRow(`SELECT exists(SELECT 1 FROM book where title=$1)`, name).Scan(&checker)
	if err != nil {
		return err
	}
	if checker {
		return errors.New("name constraint. You are trying to save this name 2nd time")
	}

	_, err = r.db.Exec("INSERT INTO book (title) VALUES ($1)", name)
	if err != nil {
		log.Printf("Error creating book: %v", err)
		return err
	}
	return nil
}

func (r *PostgresBookRepository) DeleteBook(id string) error {
	bookId, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM book WHERE id = $1", bookId)
	if err != nil {
		log.Printf("Error deleting book: %v", err)
		return err
	}
	return nil
}

func (r *PostgresBookRepository) GetAllBooks() ([]models.Book, error) {
	rows, err := r.db.Query("SELECT id, title FROM book")
	if err != nil {
		log.Printf("Error retrieving books: %v", err)
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title); err != nil {
			log.Printf("Error scanning book row: %v", err)
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *PostgresBookRepository) CreateAnswer(bookId string, sectionType string, answer []string) error {
	id, err := strconv.Atoi(bookId)
	if err != nil {
		return err
	}
	var checker bool
	err = r.db.QueryRow(`SELECT exists(SELECT 1 FROM book WHERE id=$1)`, id).Scan(&checker)
	if err != nil {
		return errors.New("error while checking if book exists")
	}
	if !checker {
		return errors.New("book not found")
	}

	if sectionType != "READING" && sectionType != "LISTENING" {
		return errors.New("invalid section type")
	}

	_, err = r.db.Exec(`INSERT INTO answer (book_id, section_type, section_answer) VALUES ($1, $2, $3)`,
		id, sectionType, pq.Array(answer))
	if err != nil {
		return errors.New("error while inserting answer into the database")
	}

	return nil
}

func (r *PostgresBookRepository) DeleteAnswer(answerId string) error {
	id, err := strconv.Atoi(answerId)
	if err != nil {
		return err
	}
	var exists bool
	err = r.db.QueryRow(`SELECT exists(SELECT 1 FROM answer WHERE id=$1)`, id).Scan(&exists)
	if err != nil {
		return errors.New("error while checking if answer exists")
	}
	if !exists {
		return errors.New("answer not found")
	}
	_, err = r.db.Exec(`DELETE FROM answer WHERE id = $1`, id)
	if err != nil {
		return errors.New("error while deleting answer from the database")
	}

	return nil
}

func (r *PostgresBookRepository) GetAnswerByBookId(bookId string) ([]models.Answer, error) {
	id, err := strconv.Atoi(bookId)
	if err != nil {
		return nil, err
	}

	var exists bool
	err = r.db.QueryRow(`SELECT exists(SELECT 1 FROM book WHERE id=$1)`, id).Scan(&exists)
	if err != nil {
		return nil, errors.New("error while checking if book exists")
	}
	if !exists {
		return nil, errors.New("book not found")
	}

	rows, err := r.db.Query(`SELECT id, book_id, section_type, section_answer FROM answer WHERE book_id=$1`, id)
	if err != nil {
		return nil, errors.New("error while retrieving answers from the database")
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		err := rows.Scan(&answer.ID, &answer.BookId, &answer.SectionType, pq.Array(&answer.Answer))
		if err != nil {
			return nil, errors.New("error while scanning answer row")
		}
		answers = append(answers, answer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return answers, nil
}
