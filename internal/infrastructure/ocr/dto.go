package ocr

import (
	"mime/multipart"
)

type OCRRequest struct {
	File *multipart.FileHeader
}

type OCRBulkRequest struct {
	Files []*multipart.FileHeader
}

type OCRResponse struct {
	FileName   string `json:"file_name"`
	TotalPages int    `json:"total_pages"`
	FullText   string `json:"full_text"`
	Pages      []Page `json:"pages"`
}

type Page struct {
	PageNumber int    `json:"page_number"`
	Text       string `json:"text"`
	Words      []Word `json:"words"`
}

type Word struct {
	Text        string   `json:"text"`
	Confidence  float64  `json:"confidence"`
	BoundingBox [][2]int `json:"bounding_box"`
}

type CORResponse struct {
	FileName          string `json:"file_name"`
	LastName          string `json:"last_name"`
	FullName          string `json:"full_name"`
	StudentNumber     string `json:"student_number"`
	StartAcademicYear string `json:"start_academic_year"`
	EndAcademicYear   string `json:"end_academic_year"`
	Term              int    `json:"term"`
	ProgramDesc       string `json:"program_desc"`
	ProgramCode       string `json:"program_code"`
	YearLevel         int    `json:"year_level"`
	Campus            string `json:"campus"`
	Section           int    `json:"section"`
}
