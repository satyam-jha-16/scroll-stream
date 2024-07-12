package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/satyam-jha-16/streamlite/backend/initializers"
	"github.com/satyam-jha-16/streamlite/backend/models"
)

func init(){
  if err := godotenv.Load(); err != nil{
    log.Fatal("error loading the env file")
  }
 
}

func main(){
  initializers.DB.AutoMigrate(&models.User{})
  initializers.DB.AutoMigrate(&models.Video{})
}
