package model

type Image struct {
	ID         string `gorm:"primary_key"`
	Filename   string
	Original   []byte
	Quality100 []byte
	Quality75  []byte
	Quality50  []byte
	Quality25  []byte
}
