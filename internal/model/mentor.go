package model

type Mentor struct {
	ID          string `json:"id" gorm:"primaryKey"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Avatar      string `json:"avatar"`
	CareerLever string `json:"career_level"`
	Disabled    bool   `json:"-" gorm:"disabled"`
	CreatedAt   int64  `json:"-" gorm:"created_at"`
}

type CreateMentorRequest struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Avatar      string `json:"avatar"`
}

type UpdateMentorRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Avatar      string `json:"avatar"`
}
