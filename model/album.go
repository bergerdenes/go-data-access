package model

type Album struct {
	ID     int64   `json:"id" yaml:"id"`
	Title  string  `json:"title" yaml:"title"`
	Artist string  `json:"artist" yaml:"artist"`
	Price  float32 `json:"price" yaml:"price"`
}

type Created struct {
	ID int64 `json:"id" yaml:"id"`
}
