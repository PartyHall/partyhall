package models

type Song struct {
	Id       int    `json:"id" db:"id"`
	Filename string `json:"filename" db:"filename"`
	Artist   string `json:"artist" db:"artist"`
	Title    string `json:"title" db:"title"`
	Format   string `json:"format" db:"format"`
}
