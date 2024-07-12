package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init(){
  if err := godotenv.Load(); err != nil{
    log.Fatal("error loading the env file")
  }
  
}

func setUpRoutes() *gin.Engine {
  
  router := gin.Default()

  config := cors.DefaultConfig()
  config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173", "*"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

  return router
}


func main(){
  router := setUpRoutes()
  
  port := os.Getenv("SERVER_PORT")
  if port == "" {
    port = "8000"
  }

  router.Run(":" + port)
}
