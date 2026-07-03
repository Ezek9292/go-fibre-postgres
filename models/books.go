package models

type Book struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}
