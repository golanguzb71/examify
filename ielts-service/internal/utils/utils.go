package utils

import (
	"database/sql"
	"fmt"
)

func OffSetGenerator(page, size *int32) int {
	if page == nil || *page < 1 {
		p := int32(1)
		page = &p
	}
	if size == nil || *size < 1 {
		s := int32(10)
		size = &s
	}

	return int(*size * (*page - 1))
}

func CalculateBandScore(correctCount int) float64 {
	switch {
	case correctCount >= 39:
		return 9.0
	case correctCount >= 37:
		return 8.5
	case correctCount >= 35:
		return 8.0
	case correctCount >= 32:
		return 7.5
	case correctCount >= 30:
		return 7.0
	case correctCount >= 26:
		return 6.5
	case correctCount >= 23:
		return 6.0
	case correctCount >= 18:
		return 5.5
	case correctCount >= 16:
		return 5.0
	case correctCount >= 13:
		return 4.5
	case correctCount >= 10:
		return 4.0
	case correctCount >= 6:
		return 3.5
	case correctCount >= 4:
		return 3.0
	case correctCount >= 1:
		return 2.5
	default:
		return 2.0
	}
}

func UpdateOverallScore(examID string, db *sql.DB) error {
	query := `
        WITH scores AS (
            SELECT 
                COALESCE((SELECT band_score FROM reading_detail WHERE exam_id = $1), 0) as reading_score,
                COALESCE((SELECT band_score FROM listening_detail WHERE exam_id = $1), 0) as listening_score,
                COALESCE((SELECT AVG(task_band_score) FROM writing_detail WHERE exam_id = $1), 0) as writing_score,
                COALESCE((SELECT part_band_score FROM speaking_detail WHERE exam_id = $1), 0) as speaking_score
        )
        UPDATE exam
        SET over_all_band_score = (
            SELECT 
                CASE 
                    WHEN remainder >= 0.75 THEN whole + 1
                    WHEN remainder >= 0.25 THEN whole + 0.5
                    ELSE whole
                END
            FROM (
                SELECT 
                    FLOOR((reading_score + listening_score + writing_score + speaking_score) / 4.0) as whole,
                    ((reading_score + listening_score + writing_score + speaking_score) / 4.0) - FLOOR((reading_score + listening_score + writing_score + speaking_score) / 4.0) as remainder
                FROM scores
            ) subquery
        )
        FROM scores
        WHERE exam.id = $1
    `

	_, err := db.Exec(query, examID)
	if err != nil {
		return fmt.Errorf("failed to update overall score for exam %s: %w", examID, err)
	}

	return nil
}
