package service

import (
	"example/data-access/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type WebService struct {
	repo repository.AlbumRepository
}

func NewWebService(repo *repository.AlbumRepository) *WebService {
	return &WebService{repo: *repo}
}

// GetAlbums responds with the list of all albums as JSON.
func (w *WebService) GetAlbums(c *gin.Context) {
	albums, err := w.repo.Albums()
	if err != nil {
		log.Printf("Unable to get Album list from DB")
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func (w *WebService) AddAlbum(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		log.Fatalf("Unable to parse input data as Album")
		return
	}

	// Add the new album to DB
	newId, err := w.repo.AddAlbum(newAlbum)
	if err != nil {
		log.Printf("Unable to add new Album to DB: %v", newAlbum)
		return
	}
	created := Created{ID: newId}
	c.IndentedJSON(http.StatusCreated, created)

}

// GetAlbumById locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func (w *WebService) GetAlbumById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Unable to convert ID to number: %s", idStr)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album not found with ID: %s", idStr)})
		return
	}

	album, err := w.repo.AlbumByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album not found with ID: %s", idStr)})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

// DeleteAlbumById deletes the album matching ID
func (w *WebService) DeleteAlbumById(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Unable to convert ID to number: %s", idStr)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album not found with ID: %s", idStr)})
		return
	}

	err = w.repo.DeleteAlbumByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("album not found with ID: %d", id)})
		return
	}

	c.Status(http.StatusNoContent)
}
