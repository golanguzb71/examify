package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"ielts-service/proto/pb"
	"log"
	"math"
)

func setExtraFieldResult(result *pb.Result, db *sql.DB) {
	examId := result.ExamId

	err := db.QueryRow(`SELECT band_score FROM listening_detail WHERE exam_id = $1`, examId).Scan(&result.Listening)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			result.Listening = "N/A"
		} else {
			log.Println("Error fetching listening band score:", err.Error())
		}
	}

	err = db.QueryRow(`SELECT band_score FROM reading_detail WHERE exam_id = $1`, examId).Scan(&result.Reading)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			result.Reading = "N/A"
		} else {
			log.Println("Error fetching reading band score:", err.Error())
		}
	}

	var task1Score, task2Score float64
	err = db.QueryRow(`SELECT task_band_score FROM writing_detail WHERE exam_id = $1 AND task_number = 1`, examId).Scan(&task1Score)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error fetching writing task 1 score:", err.Error())
	}

	err = db.QueryRow(`SELECT task_band_score FROM writing_detail WHERE exam_id = $1 AND task_number = 2`, examId).Scan(&task2Score)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error fetching writing task 2 score:", err.Error())
	}

	if task1Score != 0 || task2Score != 0 {
		result.Writing = roundToNearestHalf((task1Score + 2*task2Score) / 3)
	} else {
		result.Writing = "N/A"
	}

	var part1Score, part2Score, part3Score float64
	err = db.QueryRow(`SELECT part_band_score FROM speaking_detail WHERE exam_id = $1 AND part_number = 1`, examId).Scan(&part1Score)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error fetching speaking part 1 score:", err.Error())
	}

	err = db.QueryRow(`SELECT part_band_score FROM speaking_detail WHERE exam_id = $1 AND part_number = 2`, examId).Scan(&part2Score)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error fetching speaking part 2 score:", err.Error())
	}

	err = db.QueryRow(`SELECT part_band_score FROM speaking_detail WHERE exam_id = $1 AND part_number = 3`, examId).Scan(&part3Score)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println("Error fetching speaking part 3 score:", err.Error())
	}

	if part1Score != 0 || part2Score != 0 || part3Score != 0 {
		result.Speaking = roundToNearestHalf((part1Score + part2Score + part3Score) / 3)
	} else {
		result.Speaking = "N/A"
	}
	
}

func roundToNearestHalf(score float64) string {
	return fmt.Sprintf("%.1f", math.Round(score*2)/2)
}
