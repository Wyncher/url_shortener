package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"url_shortener/config"

	"url_shortener/services"
)

type Handler struct{}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/")
	v1.GET("/shrink", services.Shrink) //принимает link и возвращает short_link
	v1.GET("/", services.Redirect)     //принимает short_link и делает редирект
	r.Run(fmt.Sprintf(":%v", config.SERVER_PORT))
}
