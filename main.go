package main

import (
	"github.com/gin-gonic/gin"
	"jwt-psql/handlers"
	"jwt-psql/models"
	"log"
)

func main() {
	r := gin.Default()

	models.Init()
	public := r.Group("/api")

	public.POST("/register", handlers.Register)

	err := r.Run(":8080")

	if err != nil {
		log.Fatal(err)
		return
	}
}
