package models

type Livro struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Titulo string `json:"titulo" binding:"required"`
	Autor  string `json:"autor" binding:"required"`
}
