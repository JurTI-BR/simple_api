package models

type Livro struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Titulo    string `json:"titulo"`
	Autor     string `json:"autor"`
	Ano       int    `json:"ano"`
	ImagePath string `json:"image_path"`
}
