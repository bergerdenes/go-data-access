package config

import "fmt"

func GetDbUrl() string {
	user := "postgres"
	password := "postgres"
	host := "localhost"
	port := 5432
	db := "album"
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, db)
}
