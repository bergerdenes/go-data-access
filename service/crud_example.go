package service

import (
	"example/data-access/repository"
	"fmt"
	"log"
)

type CrudService struct {
	repo repository.AlbumRepository
}

//func NewCrudService(repo *repository.AlbumRepository) *CrudService {
//	return &CrudService{repo: *repo}
//}

func (s *CrudService) SingleRowQuery() {
	alb, err := s.repo.AlbumByID(4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)
}

func (s *CrudService) MultiRowQuery() {
	albums, err := s.repo.AlbumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
}

func (s *CrudService) InsertRow() {
	albID, err := s.repo.AddAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}
