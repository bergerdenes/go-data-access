package repository

import (
	"context"
	"errors"
	"example/data-access/model"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlbumRepository struct {
	connection *pgxpool.Pool
}

type Album = model.Album

func New(connection *pgxpool.Pool) *AlbumRepository {
	return &AlbumRepository{connection: connection}
}

// Albums queries for all albums in database.
func (r *AlbumRepository) Albums() ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := r.connection.Query(context.Background(), "SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("albums: %v", err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albums: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albums: %v", err)
	}
	return albums, nil
}

// AlbumsByArtist queries for albums that have the specified artist name.
func (r *AlbumRepository) AlbumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := r.connection.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// AlbumByID queries for the album with the specified ID.
func (r *AlbumRepository) AlbumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := r.connection.QueryRow(context.Background(), "SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// AddAlbum adds the specified album to the database,
// returning the album ID of the new entry
func (r *AlbumRepository) AddAlbum(alb Album) (int64, error) {
	var newId int64
	err := r.connection.QueryRow(
		context.Background(),
		"INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
		alb.Title, alb.Artist, alb.Price,
	).Scan(&newId)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return newId, nil
}

func (r *AlbumRepository) DeleteAlbumByID(id int64) error {
	result, err := r.connection.Exec(context.Background(), "DELETE FROM album WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleteAlbum: %v", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("deleteAlbum: album with ID %d was not found", id)
	}
	return nil
}
