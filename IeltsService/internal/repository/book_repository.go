package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	client "ielts-service/internal/grpc_clients"
	"ielts-service/internal/models"
	"ielts-service/internal/utils"
	"ielts-service/proto/pb"
	"log"
	"strconv"
)

type PostgresRepository struct {
	db         *sql.DB
	userClient *client.UserClient
}

func NewPostgresRepository(db *sql.DB, userClient *client.UserClient) *PostgresRepository {
	return &PostgresRepository{db: db, userClient: userClient}
}

func (r *PostgresRepository) CreateBook(name string) error {
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

func (r *PostgresRepository) DeleteBook(id string) error {
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

func (r *PostgresRepository) GetAllBooks() ([]models.Book, error) {
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

func (r *PostgresRepository) CreateAnswer(bookId string, sectionType string, answer []string) error {
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

	if sectionType != "READING" && sectionType != "LISTENING" && sectionType != "WRITING" {
		return errors.New("invalid section type")
	}

	_, err = r.db.Exec(`INSERT INTO answer (book_id, section_type, section_answer) VALUES ($1, $2, $3)`,
		id, sectionType, pq.Array(answer))
	if err != nil {
		log.Println(err)
		return errors.New("error while inserting answer into the database")
	}

	return nil
}

func (r *PostgresRepository) DeleteAnswer(answerId string) error {
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

func (r *PostgresRepository) GetAnswerByBookId(bookId string) ([]models.Answer, error) {
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

func (r *PostgresRepository) CreateExam(userID, bookID int32) (*string, error) {
	var id string
	err := r.db.QueryRow(
		`INSERT INTO exam(id, user_id, book_id) VALUES ($1, $2, $3) RETURNING id`,
		uuid.New().String(), userID, bookID,
	).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (r *PostgresRepository) GetExamsByUserId(userID, page, size int32) (*pb.GetExamByUserIdResponse, error) {
	r.db.Query(`SELECT id, user_id, book_id, over_all_band_score , created_at FROM exam where user_id=$1 order by created_at desc  offset $2 limit $3`)
	return nil, nil
}

func (r *PostgresRepository) GetTopExamResults(dataframe string, page, size int32) (*pb.GetTopExamResult, error) {
	baseQuery := `
		SELECT e.id, e.user_id, b.title, e.over_all_band_score, b.created_at
		FROM exam e
		JOIN book b ON e.book_id = b.id 
		WHERE `

	var timeframeCondition string
	switch dataframe {
	case "MONTHLY":
		timeframeCondition = `e.created_at >= date_trunc('month', CURRENT_DATE)`
	case "DAILY":
		timeframeCondition = `e.created_at >= CURRENT_DATE`
	case "WEEKLY":
		timeframeCondition = `e.created_at >= date_trunc('week', CURRENT_DATE)`
	default:
		return nil, fmt.Errorf("invalid dataframe: %s", dataframe)
	}

	finalQuery := baseQuery + timeframeCondition + `
		ORDER BY e.over_all_band_score
		LIMIT $1 OFFSET $2`

	offset := utils.OffSetGenerator(&page, &size)
	rows, err := r.db.Query(finalQuery, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*pb.Result
	for rows.Next() {
		var result pb.Result
		var userId string
		err = rows.Scan(&result.ExamId, &userId, &result.BookName, &result.Overall, &result.CreatedAt)
		if err != nil {
			return nil, err
		}

		user := r.userClient.GetUserByPhoneNumberOrChatIdOrId(nil, nil, &userId)
		result.User = user

		setExtraFieldResult(&result, r.db)
		fmt.Println(&result)
		results = append(results, &result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetTopExamResult{Results: results}, nil
}

func (r *PostgresRepository) CreateAttemptInline(examID string, userAnswer []string, sectionType string) error {
	return nil
}

func (r *PostgresRepository) CreateAttemptOutline(examID string, answers *pb.CreateOutlineAttemptRequest) error {
	return nil
}
