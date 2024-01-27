package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/data-access/config"
	"example/data-access/repository"
	"example/data-access/service"
)

// var crudSvc *service.CrudService

func main() {
	//gin.SetMode(gin.ReleaseMode)
	setupLogging()
	connection := connect()
	repo := repository.New(connection)
	defer connection.Close()
	connectionTest(connection)

	webSvc := service.NewWebService(repo)
	router := gin.Default()
	router.GET("/albums", webSvc.GetAlbums)
	router.POST("/album", webSvc.AddAlbum)
	router.GET("/album/:id", webSvc.GetAlbumById)
	router.DELETE("/album/:id", webSvc.DeleteAlbumById)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Unable to start web service")
	}

	//crudSvc = service.NewCrudService(repo)
	//tryCrudFunctions()

}

//func tryCrudFunctions() {
//	crudSvc.MultiRowQuery()
//	crudSvc.SingleRowQuery()
//	crudSvc.InsertRow()
//}

func setupLogging() {
	log.SetPrefix("data-access: ")
	log.SetFlags(0)
}

func connect() *pgxpool.Pool {
	poolConfig, errParse := pgxpool.ParseConfig(config.GetDbUrl())
	if errParse != nil {
		log.Fatal("Unable to parse DB URL")

	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	return conn
}

func connectionTest(connection *pgxpool.Pool) {
	errPing := connection.Ping(context.Background())
	if errPing != nil {
		log.Fatal(errPing)
	}
	log.Println("Connected!")
}
