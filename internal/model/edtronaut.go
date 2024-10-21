package model

type EdtronautAPIResponse[T any] struct {
	Data    T   `json:"data"`
	Success int `json:"success"`
}

type EdtronautUser struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	ID        string `json:"id"`
	SurName   string `json:"sur_name"`
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
}

type EdtronautCourse struct {
	AllowBooking bool   `json:"allow_booking"`
	Id           string `json:"id"`
	CourseName   string `json:"course_name"`
	CourseType   int    `json:"start_date"`
}

type EdtronautPagination struct {
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Count     int `json:"count"`
}

type EdtronautPaginationResponse[T any] struct {
	Data struct {
		Courses []T `json:"courses"`
	} `json:"data"`
	Pagination EdtronautPagination `json:"pagination"`
	Success    int                 `json:"success"`
}
