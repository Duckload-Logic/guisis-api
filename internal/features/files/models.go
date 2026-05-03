package files

import "time"

// File represents a file asset, used for both business logic and data persistence.
type File struct {
	ID        string    `db:"id"         json:"id"`
	FileName  string    `db:"file_name"  json:"fileName"`
	FileURL   string    `db:"file_url"   json:"fileUrl"`
	FileType  string    `db:"file_type"  json:"fileType"`
	FileSize  int64     `db:"file_size"  json:"fileSize"`
	MimeType  string    `db:"mime_type"  json:"mimeType"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	DeletedAt time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// OCRResult represents the result of an OCR process.
type OCRResult struct {
	FileID         string    `db:"file_id"         json:"fileId"`
	RawText        string    `db:"raw_text"        json:"rawText"`
	StructuredData string    `db:"structured_data" json:"structuredData"`
	EngineV        string    `db:"engine_v"        json:"engineV"`
	Confidence     float64   `db:"confidence"      json:"confidence"`
	CreatedAt      time.Time `db:"created_at"      json:"createdAt"`
}
