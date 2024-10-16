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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	//if err != nil || count >= 2 {
	//	return nil, errors.New("you can create exam 2 times in a day")
	//}

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
	offset := (page - 1) * size
	query := `
        WITH exam_data AS (
            SELECT 
                e.id AS exam_id, b.title AS book_name, e.created_at, 
                e.over_all_band_score, e.status,
                COALESCE(AVG(sd.part_band_score), 0) AS speaking_score,
                COALESCE(AVG(wd.task_band_score), 0) AS writing_score,
                COALESCE(ld.band_score, 0) AS listening_score,
                COALESCE(rd.band_score, 0) AS reading_score
            FROM exam e
            JOIN book b ON e.book_id = b.id
            LEFT JOIN speaking_detail sd ON e.id = sd.exam_id
            LEFT JOIN writing_detail wd ON e.id = wd.exam_id
            LEFT JOIN listening_detail ld ON e.id = ld.exam_id
            LEFT JOIN reading_detail rd ON e.id = rd.exam_id
            WHERE e.user_id = $1
            GROUP BY e.id, b.title, e.created_at, e.over_all_band_score, e.status , ld.band_score , rd.band_score
        )
        SELECT exam_id, book_name, created_at, over_all_band_score, status,
               speaking_score, writing_score, listening_score, reading_score,
               COUNT(*) OVER() AS total_count
        FROM exam_data
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.Query(query, userID, size, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*pb.GetExamAbsResult
	var totalCount int32

	for rows.Next() {
		var examId, bookName, status string
		var createdAt time.Time
		var overallScore, speakingScore, writingScore, listeningScore, readingScore float64

		err := rows.Scan(&examId, &bookName, &createdAt, &overallScore, &status,
			&speakingScore, &writingScore, &listeningScore, &readingScore,
			&totalCount)
		if err != nil {
			return nil, err
		}

		remainTime := int32(0)
		var remainSection string
		if status == "PENDING" {
			endTime := createdAt.Add(4 * time.Hour)
			if remain := time.Until(endTime); remain > 0 {
				remainTime = int32(remain.Seconds())
			}
			checkIsHaveRemainSection(&remainSection, r.db, examId)
		}

		results = append(results, &pb.GetExamAbsResult{
			ExamId:               examId,
			BookName:             bookName,
			CreatedAt:            createdAt.Format(time.RFC3339),
			Overall:              fmt.Sprintf("%.1f", overallScore),
			Speaking:             fmt.Sprintf("%.1f", speakingScore),
			Writing:              fmt.Sprintf("%.1f", writingScore),
			Listening:            fmt.Sprintf("%.1f", listeningScore),
			Reading:              fmt.Sprintf("%.1f", readingScore),
			Status:               status,
			RemainTimeForEndExam: remainTime,
			RemainSection:        remainSection,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := int32(math.Ceil(float64(totalCount) / float64(size)))
	return &pb.GetExamByUserIdResponse{
		Results:        results,
		TotalPageCount: totalPages,
	}, nil
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
		timeframeCondition = `e.created_at >= CURRENT_DATE - INTERVAL '7 days'`
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
	parsedUUID, err := uuid.Parse(examID)
	if err != nil {
		return err
	}
	var checker bool
	err = r.db.QueryRow(`SELECT exists(SELECT 1 FROM exam where id=$1 and status='PENDING')`, parsedUUID).Scan(&checker)
	if err != nil || !checker {
		return errors.New("exam not found or Exam already finished")
	}
	var bookID int
	query := `SELECT book_id FROM exam WHERE id = $1`
	err = r.db.QueryRow(query, examID).Scan(&bookID)
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
	err = r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM exam where id=$1 and status='PENDING')`, parsedUUID).Scan(&checker)
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
	err = utils.UpdateOverallScore(id, r.db)
	if err != nil {
		return err
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

func (r *PostgresRepository) CreateAttemptOutlineSpeaking(req *pb.CreateOutlineAttemptRequestSpeaking) error {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM exam WHERE status = 'PENDING' AND id = $1
		)
	`, req.ExamId).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking for pending exam: %v", err)
	}
	if !exists {
		return fmt.Errorf("no pending speaking attempts for exam ID: %s", req.ExamId)
	}

	fileName := fmt.Sprintf("%s_part%d_%s.wav", req.ExamId, req.PartNumber, uuid.New().String())
	filePath := filepath.Join("voice_answers", fileName)

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	err = os.WriteFile(filePath, req.VoiceAnswer, 0644)
	if err != nil {
		return fmt.Errorf("error writing voice answer to file: %v", err)
	}

	voiceURL := fmt.Sprintf("/voice_answers/%s", fileName)

	resp, err := r.integrationClient.GetResultSpeakingPart(req)
	if err != nil {
		return fmt.Errorf("error getting result from integration client: %v", err)
	}

	transcription := map[string]string{
		"question":      resp.Transcription.Question,
		"feedback":      resp.Transcription.Feedback,
		"transcription": resp.Transcription.Transcription,
	}

	transcriptionJSON, err := json.Marshal(transcription)
	if err != nil {
		return fmt.Errorf("error marshaling transcription to JSON: %v", err)
	}

	if err := validateIELTSBandScores(
		resp.FluencyScore,
		resp.GrammarScore,
		resp.VocabularyScore,
		resp.CoherenceScore,
		resp.TopicDevScore,
		resp.RelevanceScore,
		resp.PartBandScore,
	); err != nil {
		return err
	}

	_, err = r.db.Exec(`
	INSERT INTO speaking_detail (
		id, exam_id, part_number, fluency_score, grammar_score, vocabulary_score,
		coherence_score, topic_dev_score, relevance_score, transcription, voice_url, part_band_score, word_count
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, ARRAY[$11], $12, $13
	)
	ON CONFLICT (exam_id, part_number) DO UPDATE SET
		fluency_score = (speaking_detail.fluency_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.fluency_score) / array_length(speaking_detail.voice_url, 1),
		grammar_score = (speaking_detail.grammar_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.grammar_score) / array_length(speaking_detail.voice_url, 1),
		vocabulary_score = (speaking_detail.vocabulary_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.vocabulary_score) / array_length(speaking_detail.voice_url, 1),
		coherence_score = (speaking_detail.coherence_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.coherence_score) / array_length(speaking_detail.voice_url, 1),
		topic_dev_score = (speaking_detail.topic_dev_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.topic_dev_score) / array_length(speaking_detail.voice_url, 1),
		relevance_score = (speaking_detail.relevance_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.relevance_score) / array_length(speaking_detail.voice_url, 1),
		part_band_score = (speaking_detail.part_band_score * (array_length(speaking_detail.voice_url, 1) - 1) + EXCLUDED.part_band_score) / array_length(speaking_detail.voice_url, 1),
		word_count = speaking_detail.word_count + EXCLUDED.word_count,
		transcription = speaking_detail.transcription || jsonb_build_array(EXCLUDED.transcription),
		voice_url = array_append(speaking_detail.voice_url, EXCLUDED.voice_url[1])
`,
		uuid.New(), req.ExamId, req.PartNumber, resp.FluencyScore, resp.GrammarScore, resp.VocabularyScore,
		resp.CoherenceScore, resp.TopicDevScore, resp.RelevanceScore, transcriptionJSON, voiceURL, resp.PartBandScore, resp.WordCount,
	)
	if err != nil {
		return fmt.Errorf("error executing insert/update query: %v", err)
	}

	err = utils.UpdateOverallScore(req.ExamId, r.db)
	if err != nil {
		return fmt.Errorf("error updating overall score: %v", err)
	}

	return nil
}

func (r *PostgresRepository) GetResultsInlineBySection(section string, examId string) (*pb.GetResultResponse, error) {
	query := "SELECT id, band_score, user_answer, created_at FROM "
	switch section {
	case "LISTENING":
		query += "listening_detail"
	case "READING":
		query += "reading_detail"
	default:
		return nil, errors.New("invalid section format. Only LISTENING, READING approved")
	}
	query += " WHERE exam_id=$1"

	var result pb.GetResultResponse
	var userAnswerRaw []byte

	err := r.db.QueryRow(query, examId).Scan(&result.Id, &result.BandScore, &userAnswerRaw, &result.CreatedAt)
	if err != nil {
		return nil, err
	}
	var userAnswers []struct {
		IsTrue     bool   `json:"is_true"`
		TrueAnswer string `json:"true_answer"`
		UserAnswer string `json:"user_answer"`
	}
	if err = json.Unmarshal(userAnswerRaw, &userAnswers); err != nil {
		return nil, err
	}
	result.Answers = make([]*pb.UserAnswer, len(userAnswers))
	for i, answer := range userAnswers {
		result.Answers[i] = &pb.UserAnswer{
			UserAnswer: answer.UserAnswer,
			TrueAnswer: answer.TrueAnswer,
			IsTrue:     answer.IsTrue,
		}
	}
	return &result, nil
}

func (r *PostgresRepository) GetResultOutlineSpeaking(req *pb.GetResultOutlineAbsRequest) (*pb.GetResultOutlineSpeakingResponse, error) {
	rows, err := r.db.Query(`
		SELECT part_number, fluency_score, grammar_score, vocabulary_score, coherence_score, 
		       topic_dev_score, relevance_score, transcription, voice_url, part_band_score, created_at 
		FROM speaking_detail 
		WHERE exam_id=$1`, req.ExamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &pb.GetResultOutlineSpeakingResponse{}
	var answers []*pb.SpeakingPartsResponse

	for rows.Next() {
		var part pb.SpeakingPartsResponse
		var transcription []byte
		var voiceUrls []string

		err = rows.Scan(
			&part.PartNumber, &part.FluencyScore, &part.GrammarScore, &part.VocabularyScore,
			&part.CoherenceScore, &part.TopicDevScore, &part.RelevanceScore, &transcription,
			pq.Array(&voiceUrls), &part.PartBandScore)
		if err != nil {
			return nil, err
		}
		var transcriptionData struct {
			Question      string `json:"question"`
			Feedback      string `json:"feedback"`
			Transcription string `json:"transcription"`
		}
		if err = json.Unmarshal(transcription, &transcriptionData); err != nil {
			return nil, err
		}

		part.Transcription = &pb.Transcription{
			Question:      transcriptionData.Question,
			Feedback:      transcriptionData.Feedback,
			Transcription: transcriptionData.Transcription,
		}

		if len(voiceUrls) > 0 {
			part.VoiceUrl = voiceUrls[0]
		}
		answers = append(answers, &part)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result.Answers = answers
	return result, nil
}

func (r *PostgresRepository) GetResultOutlineWriting(req *pb.GetResultOutlineAbsRequest) (*pb.GetResultOutlineWritingResponse, error) {
	rows, err := r.db.Query(`
		SELECT task_number, response, feedback, coherence_score, grammar_score, 
		lexical_resource_score, task_achievement_score, task_band_score, created_at 
		FROM writing_detail 
		WHERE exam_id=$1`, req.ExamId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := &pb.GetResultOutlineWritingResponse{}
	var answer []*pb.OutlineWritingResponseAbs
	var response []byte

	for rows.Next() {
		var ans pb.OutlineWritingResponseAbs
		err = rows.Scan(&ans.TaskNumber, &response, &ans.Feedback, &ans.CoherenceScore, &ans.GrammarScore,
			&ans.LexicalResourceScore, &ans.TaskAchievementScore, &ans.TaskBandScore, &ans.CreatedAt)
		if err != nil {
			return nil, err
		}

		var userResponse struct {
			Question   string `json:"question"`
			UserAnswer string `json:"user_answer"`
		}

		if err = json.Unmarshal(response, &userResponse); err != nil {
			return nil, err
		}

		ans.Question = userResponse.Question
		ans.UserAnswer = userResponse.UserAnswer

		answer = append(answer, &ans)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	result.Answers = answer
	return result, nil
}

func validateIELTSBandScores(scores ...float32) error {
	validScores := map[float32]bool{
		1:   true,
		1.5: true,
		2:   true,
		2.5: true,
		3:   true,
		3.5: true,
		4:   true,
		4.5: true,
		5:   true,
		5.5: true,
		6:   true,
		6.5: true,
		7:   true,
		7.5: true,
		8:   true,
		8.5: true,
	}

	for _, score := range scores {
		if !validScores[score] {
			return fmt.Errorf("invalid score: %.2f; valid scores are 1, 1.5, ..., 8, 8.5", score)
		}
	}
	return nil
}

func checkIsHaveRemainSection(remainSection *string, db *sql.DB, examId string) {
	var checkerListening, checkerReading, checkerSpeaking, checkerWriting bool
	checkerListeningQuery := `SELECT exists(SELECT 1 FROM listening_detail where exam_id=$1)`
	_ = db.QueryRow(checkerListeningQuery, examId).Scan(&checkerListening)
	if !checkerListening {
		*remainSection = "LISTENING"
		return
	}
	checkerReadingQuery := `SELECT exists(SELECT 1 FROM reading_detail where exam_id=$1)`
	_ = db.QueryRow(checkerReadingQuery, examId).Scan(&checkerReading)
	if !checkerReading {
		*remainSection = "READING"
		return
	}
	checkerSpeakingQuery := `SELECT exists(SELECT 1 FROM speaking_detail where exam_id=$1)`
	_ = db.QueryRow(checkerSpeakingQuery, examId).Scan(&checkerSpeaking)
	if !checkerSpeaking {
		*remainSection = "SPEAKING"
		return
	}
	checkerWritingQuery := `SELECT exists(SELECT 1 FROM writing_detail where exam_id=$1)`
	_ = db.QueryRow(checkerWritingQuery, examId).Scan(&checkerWriting)
	if !checkerWriting {
		*remainSection = "WRITING"
		return
	}
}
