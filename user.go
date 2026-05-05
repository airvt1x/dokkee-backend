package dokkee

type User struct {
	Id         int     `json:"id" db:"id"`
	Username   string  `json:"username" binding:"required"`
	Password   string  `json:"password" binding:"required"`
	FirstName  string  `json:"first_name" binding:"required"`
	LastName   string  `json:"last_name" binding:"required"`
	MiddleName string  `json:"middle_name"`
	Email      string  `json:"email" binding:"required"`
	Phone      string  `json:"phone" binding:"required"`
	Balance    float64 `json:"balance"`
}

type UpdateProfileInput struct {
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	MiddleName *string `json:"middle_name"`
	Phone      *string `json:"phone"`
}
