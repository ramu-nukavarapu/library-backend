package main

// Custom type for Book Status
type Status string

const (
	Read    Status = "read"
	Reading Status = "reading"
	ToRead  Status = "to_read"
)

// user struct to store user details
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"-"`
}

// book struct to store book details
type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title"`
	Status Status `json:"status" gorm:"default:to_read"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	UserID int    `json:"user_id"`
}
