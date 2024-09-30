package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	client "ielts-service/internal/grpc_clients"
	"ielts-service/internal/models"
	"ielts-service/internal/utils"
	"ielts-service/proto/pb"
	"log"
	"math"
	"strconv"
	"strings"
)

type PostgresRepository struct {
	db                *sql.DB
	userClient        *client.UserClient
	integrationClient *client.IntegrationClient
}

func NewPostgresRepository(db *sql.DB, userClient *client.UserClient, integrationClient *client.IntegrationClient) *PostgresRepository {
	return &PostgresRepository{db: db, userClient: userClient, integrationClient: integrationClient}
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
	var count int
	err := r.db.QueryRow(`SELECT count(*) 
                      FROM exam 
                      WHERE user_id = $1 
                      AND DATE(created_at) = CURRENT_DATE`, userID).Scan(&count)
	if err != nil || count >= 2 {
		return nil, errors.New("you can create exam 2 times in a day")
	}

	var id string
	err = r.db.QueryRow(
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
		WHERE e.status='FINISHED' and `

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

	countQuery := `
		SELECT COUNT(*) 
		FROM exam e 
		JOIN book b ON e.book_id = b.id 
		WHERE e.status='FINISHED' and ` + timeframeCondition

	var totalCount int32
	err := r.db.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	finalQuery := baseQuery + timeframeCondition + `
		ORDER BY e.over_all_band_score
		LIMIT $1 OFFSET $2`

	totalPageCount := int32(math.Ceil(float64(totalCount) / float64(size)))
	if page > totalPageCount {
		return &pb.GetTopExamResult{Results: []*pb.Result{}, TotalPageCount: totalPageCount}, nil
	}

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
		results = append(results, &result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetTopExamResult{Results: results, TotalPageCount: totalPageCount}, nil
}

func (r *PostgresRepository) CreateAttemptInline(examID string, userAnswer []string, sectionType string) error {
	var bookID int
	query := `SELECT book_id FROM exam WHERE id = $1`
	err := r.db.QueryRow(query, examID).Scan(&bookID)
	if err != nil {
		return fmt.Errorf("failed to fetch book ID for exam %s: %w", examID, err)
	}

	var correctAnswers []string
	query = `SELECT section_answer FROM answer WHERE book_id = $1 AND section_type = $2`
	err = r.db.QueryRow(query, bookID, sectionType).Scan(pq.Array(&correctAnswers))
	if err != nil {
		return fmt.Errorf("failed to fetch correct answers for book ID %d: %w", bookID, err)
	}

	if len(userAnswer) != len(correctAnswers) {
		return errors.New("number of user answers does not match the number of correct answers")
	}

	var correctCount int
	var answerDetails []models.AnswerDetail
	for i, uAnswer := range userAnswer {
		isTrue := strings.EqualFold(strings.TrimSpace(uAnswer), strings.TrimSpace(correctAnswers[i]))
		if isTrue {
			correctCount++
		}
		answerDetails = append(answerDetails, models.AnswerDetail{
			UserAnswer: uAnswer,
			TrueAnswer: correctAnswers[i],
			IsTrue:     isTrue,
		})
	}

	bandScore := utils.CalculateBandScore(correctCount)

	answerDetailsJSON, err := json.Marshal(answerDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal answer details: %w", err)
	}

	switch sectionType {
	case "READING":
		query = `
            INSERT INTO reading_detail (id, exam_id, band_score, user_answer, created_at)
            VALUES ($1, $2, $3, $4, now())`
		_, err = r.db.Exec(query, uuid.New(), examID, bandScore, answerDetailsJSON)
		if err != nil {
			return fmt.Errorf("failed to insert reading detail: %w", err)
		}
	case "LISTENING":
		query = `
            INSERT INTO listening_detail (id, exam_id, band_score, user_answer, created_at)
            VALUES ($1, $2, $3, $4, now())`
		_, err = r.db.Exec(query, uuid.New(), examID, bandScore, answerDetailsJSON)
		if err != nil {
			return fmt.Errorf("failed to insert listening detail: %w", err)
		}
	default:
		return errors.New("invalid section type: inline attempts only support READING or LISTENING")
	}
	err = utils.UpdateOverallScore(examID, r.db)
	if err != nil {
		return fmt.Errorf("failed to update overall score: %w", err)
	}
	return nil
}

func (r *PostgresRepository) CreateAttemptOutlineWriting(req *pb.CreateOutlineAttemptRequestWriting) error {
	id := req.ExamId
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	var checker = false
	err = r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM exam where id=$1 and status!='FINISHED')`, parsedUUID).Scan(&checker)
	if err != nil || !checker {
		return errors.New("exam not found or exam finished")
	}
	for i, perQua := range req.Qua {
		rpcRequest := pb.WritingTaskAbsRequest{
			Question: perQua.Question,
			Answer:   perQua.UserAnswer,
		}
		resp, err := r.integrationClient.GetResultWritingTask(&rpcRequest)
		if err != nil {
			return err
		}
		response, err := json.Marshal(perQua)
		if err != nil {
			return err
		}
		_, err = r.db.Exec(`INSERT INTO writing_detail(id, exam_id, task_number, response, feedback, coherence_score, grammar_score, lexical_resource_score, task_achievement_score, task_band_score) 
		values ($1 , $2 , $3 , $4 , $5 , $6 , $7,$8,$9,$10)`, uuid.New(), parsedUUID, i+1, response, resp.Feedback, resp.CoherenceScore, resp.GrammarScore, resp.LexicalResourceScore, resp.TaskAchievementScore, resp.TaskBandScore)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepository) UpdateBook(id string, name string) error {
	_, err := r.db.Exec(`UPDATE book SET title=$1 where id=$2`, name, id)
	if err != nil {
		return err
	}
	return nil
}
