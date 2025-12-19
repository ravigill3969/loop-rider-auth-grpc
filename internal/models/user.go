package models

type User struct {
	ID          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	FullName    string `json:"full_name" db:"full_name"`
	Password    string `json:"password" db:"password"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	BirthMonth  string `json:"birth_month" db:"birth_month"`
	BirthYear   int64  `json:"birth_year" db:"birth_year"`
	UpdatedAt   int64  `json:"updated_at" db:"updated_at"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
}
