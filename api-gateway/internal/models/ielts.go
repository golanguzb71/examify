package models

type Book struct {
	ID    int32  `json:"id"`
	Title string `json:"title"`
}

type CreateAnswer struct {
	Answers     []string `json:"answers" binding:"required"`
	SectionType string   `json:"sectionType"`
}
