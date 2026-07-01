package main

import (
	"fmt"
	"gopal-sub/auth_service/internal/database"
	"gopal-sub/auth_service/internal/user"
	"log"

	"github.com/joho/godotenv"
)

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal("env did not load")
	}

	db, err := database.New()

	if err != nil {
		log.Fatalf(`database error  wowowo  %v`, err)
	}
	newRepo := user.NewRepository(db)
	
	fmt.Println("connected")
	fmt.Println(newRepo)
	defer db.Close()
}