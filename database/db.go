package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(){
	var err error

	// Connecting string
	dsn := os.Getenv("DB_DSN")

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database: ", err)
	}

	// Connection test 
	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping database: ", err)
	}

	fmt.Println("Berhasil terhubung dengan database")
}