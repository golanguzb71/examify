package models

type Book struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
}

type Answer struct {
	ID          int32    `json:"id"`
	BookId      int32    `json:"bookId"`
	SectionType string   `json:"sectionType"`
	Answer      []string `json:"answer"`
}

type AnswerDetail struct {
	UserAnswer string `json:"user_answer"`
	TrueAnswer string `json:"true_answer"`
	IsTrue     bool   `json:"is_true"`
}
