package service_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"ielts-service/internal/repository"
	"testing"
)

func TestCreateAttemptInline(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := repository.NewPostgresRepository(db, nil)

	testCases := []struct {
		name        string
		examId      string
		userAnswer  []string
		sectionType string
		setupMock   func()
		expectError bool
	}{
		{
			name:        "Successful READING attempt",
			examId:      "123e4567-e89b-12d3-a456-426614174000",
			userAnswer:  []string{"A", "B", "C"},
			sectionType: "READING",
			setupMock: func() {
				mock.ExpectQuery("SELECT book_id FROM exam WHERE id = \\$1").
					WithArgs("123e4567-e89b-12d3-a456-426614174000").
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(1))

				mock.ExpectQuery("SELECT section_answer FROM answer WHERE book_id = \\$1 AND section_type = \\$2").
					WithArgs(1, "READING").
					WillReturnRows(sqlmock.NewRows([]string{"section_answer"}).AddRow(pq.Array([]string{"A", "B", "C"})))

				mock.ExpectExec("INSERT INTO reading_detail").
					WithArgs(sqlmock.AnyArg(), "123e4567-e89b-12d3-a456-426614174000", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name:        "Successful LISTENING attempt",
			examId:      "223e4567-e89b-12d3-a456-426614174000",
			userAnswer:  []string{"D", "E", "F"},
			sectionType: "LISTENING",
			setupMock: func() {
				mock.ExpectQuery("SELECT book_id FROM exam WHERE id = \\$1").
					WithArgs("223e4567-e89b-12d3-a456-426614174000").
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(2))

				mock.ExpectQuery("SELECT section_answer FROM answer WHERE book_id = \\$1 AND section_type = \\$2").
					WithArgs(2, "LISTENING").
					WillReturnRows(sqlmock.NewRows([]string{"section_answer"}).AddRow(pq.Array([]string{"D", "E", "F"})))

				mock.ExpectExec("INSERT INTO listening_detail").
					WithArgs(sqlmock.AnyArg(), "223e4567-e89b-12d3-a456-426614174000", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectError: false,
		},
		{
			name:        "Invalid section type",
			examId:      "323e4567-e89b-12d3-a456-426614174000",
			userAnswer:  []string{"G", "H", "I"},
			sectionType: "WRITING",
			setupMock: func() {
				mock.ExpectQuery("SELECT book_id FROM exam WHERE id = \\$1").
					WithArgs("323e4567-e89b-12d3-a456-426614174000").
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(3))

				mock.ExpectQuery("SELECT section_answer FROM answer WHERE book_id = \\$1 AND section_type = \\$2").
					WithArgs(3, "WRITING").
					WillReturnRows(sqlmock.NewRows([]string{"section_answer"}).AddRow(pq.Array([]string{"G", "H", "I"})))
			},
			expectError: true,
		},
		{
			name:        "Mismatched answer lengths",
			examId:      "423e4567-e89b-12d3-a456-426614174000",
			userAnswer:  []string{"J", "K"},
			sectionType: "READING",
			setupMock: func() {
				mock.ExpectQuery("SELECT book_id FROM exam WHERE id = \\$1").
					WithArgs("423e4567-e89b-12d3-a456-426614174000").
					WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(4))

				mock.ExpectQuery("SELECT section_answer FROM answer WHERE book_id = \\$1 AND section_type = \\$2").
					WithArgs(4, "READING").
					WillReturnRows(sqlmock.NewRows([]string{"section_answer"}).AddRow(pq.Array([]string{"J", "K", "L"})))
			},
			expectError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()
			err = repo.CreateAttemptInline(tc.examId, tc.userAnswer, tc.sectionType)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
