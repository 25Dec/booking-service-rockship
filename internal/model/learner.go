package model

// Learner -.
type Learner struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Avatar      struct {
		ContentType string `json:"content_type"`
		FileId      string `json:"file_id"`
		FileName    string `json:"file_name"`
		Key         string `json:"key"`
		PublicUrl   string `json:"public_url"`
		S3Link      string `json:"s3_link"`
	} `json:"avatar"`
}
