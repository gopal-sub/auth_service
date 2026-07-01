package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)


func New() (*sql.DB, error){
	connStr := os.Getenv("DATABASE_URL");
	db, err := sql.Open("pgx", connStr);
	if err != nil{
		return nil, err
	}
	db.SetMaxOpenConns(100);
	db.SetMaxIdleConns(50);



	err = db.Ping()
	if err != nil {
		return nil, err
	}

	stats := db.Stats()
	fmt.Printf("%v\n", stats)



	return db, nil 
}