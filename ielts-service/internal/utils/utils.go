package utils

import (
	"database/sql"
	"fmt"
	"log"
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

func MigrateUp(db *sql.DB) {
	//filePath := "./migrations/ielts_service_up.sql"
	//sqlFile, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("Error opening SQL migration file: %s", err)
	//}
	//defer sqlFile.Close()
	sqlContent := `
create table if not exists book
(
    id         serial primary key,
    title      varchar NOT NULL unique,
    created_at timestamp DEFAULT now()
);

create table if not exists answer
(
    id             serial primary key,
    book_id        int references book (id),
    section_type   varchar check ( section_type in ('READING', 'LISTENING')),
    section_answer TEXT[] NOT NULL,
    created_at     timestamp DEFAULT now(),
    UNIQUE (book_id, section_type)
);

CREATE TABLE IF NOT EXISTS exam
(
    id                  uuid PRIMARY KEY,
    user_id             int                      NOT NULL,
    book_id             int REFERENCES book (id) NOT NULL,
    over_all_band_score FLOAT CHECK (over_all_band_score >= 0 AND over_all_band_score <= 9) DEFAULT 0,
    status              varchar check ( status in ('FINISHED', 'PENDING'))                  DEFAULT 'PENDING',
    created_at          timestamp                                                           DEFAULT now(),
    end_at              timestamp                                                           DEFAULT (now() + INTERVAL '4 hours')
);


CREATE TABLE IF NOT EXISTS speaking_detail
(
    id               UUID PRIMARY KEY,
    exam_id          UUID REFERENCES exam (id) NOT NULL,
    part_number      INT CHECK (part_number IN (1, 2, 3)),
    fluency_score    FLOAT CHECK (fluency_score >= 0 AND fluency_score <= 9),
    grammar_score    FLOAT CHECK (grammar_score >= 0 AND grammar_score <= 9),
    vocabulary_score FLOAT CHECK (vocabulary_score >= 0 AND vocabulary_score <= 9),
    coherence_score  FLOAT CHECK (coherence_score >= 0 AND coherence_score <= 9),
    topic_dev_score  FLOAT CHECK (topic_dev_score >= 0 AND topic_dev_score <= 9),
    relevance_score  FLOAT CHECK (relevance_score >= 0 AND relevance_score <= 9),
    word_count       INT                       NOT NULL DEFAULT 0,
    transcription    JSONB,
    voice_url        TEXT[],
    part_band_score  FLOAT                     NOT NULL DEFAULT 0 CHECK (part_band_score >= 0 AND part_band_score <= 9),
    created_at       TIMESTAMP                          DEFAULT NOW(),
    UNIQUE (part_number, exam_id)
);

CREATE TABLE IF NOT EXISTS writing_detail
(
    id                     UUID PRIMARY KEY,
    exam_id                UUID REFERENCES exam (id)                                      NOT NULL,
    task_number            INT CHECK ( task_number IN (1, 2))                             NOT NULL,
    response               JSONB                                                          NOT NULL,
    feedback               TEXT,
    coherence_score        FLOAT CHECK (coherence_score >= 0 AND coherence_score <= 9 AND
                                        coherence_score * 2 = ROUND(coherence_score * 2)),
    grammar_score          FLOAT CHECK (grammar_score >= 0 AND grammar_score <= 9 AND
                                        grammar_score * 2 = ROUND(grammar_score * 2)),
    lexical_resource_score FLOAT CHECK (lexical_resource_score >= 0 AND lexical_resource_score <= 9 AND
                                        lexical_resource_score * 2 = ROUND(lexical_resource_score * 2)),
    task_achievement_score FLOAT CHECK (task_achievement_score >= 0 AND task_achievement_score <= 9 AND
                                        task_achievement_score * 2 = ROUND(task_achievement_score * 2)),
    task_band_score        FLOAT CHECK (task_band_score >= 0 AND task_band_score <= 9 AND
                                        task_band_score * 2 = ROUND(task_band_score * 2)) NOT NULL DEFAULT 0,
    created_at             TIMESTAMP                                                               DEFAULT NOW(),
    UNIQUE (task_number, exam_id)
);

CREATE TABLE IF NOT EXISTS listening_detail
(
    id          UUID PRIMARY KEY,
    exam_id     UUID REFERENCES exam (id)                                                   NOT NULL UNIQUE,
    band_score  FLOAT CHECK (band_score >= 0 AND band_score <= 9 AND band_score * 2 =
                                                                     ROUND(band_score * 2)) NOT NULL DEFAULT 0,
    user_answer JSONB                                                                       NOT NULL,
    created_at  TIMESTAMP                                                                            DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS reading_detail
(
    id          UUID PRIMARY KEY,
    exam_id     UUID REFERENCES exam (id)                                                   NOT NULL UNIQUE,
    band_score  FLOAT CHECK (band_score >= 0 AND band_score <= 9 AND band_score * 2 =
                                                                     ROUND(band_score * 2)) NOT NULL DEFAULT 0,
    user_answer JSONB                                                                       NOT NULL,
    created_at  TIMESTAMP                                                                            DEFAULT NOW()
);
`

	_, err := db.Exec(sqlContent)
	if err != nil {
		log.Fatalf("Error executing SQL migration: %s", err)
	}

	fmt.Println("Database migration ran successfully!")
}
