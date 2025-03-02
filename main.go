package main

import (
	"github.com/gin-gonic/gin"
	"jwt-psql/handlers"
	"jwt-psql/middlewares"
	"jwt-psql/models"
	"log"
)

func main() {
	r := gin.Default()

	models.Init()
	public := r.Group("/api")

	public.POST("/register", handlers.Register)
	public.POST("/login", handlers.Login)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", handlers.CurrentUser)

	err := r.Run(":8080")

	if err != nil {
		log.Fatal(err)
		return
	}
}
