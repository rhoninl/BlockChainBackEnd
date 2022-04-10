package main

import (
	"github.com/gin-gonic/gin"
	"main/Controller"
	"main/MiddleWare"
)

func main() {
	router := gin.Default()

	account := router.Group("/account")
	{
		account.POST("/Login", Controller.Login)
		account.POST("/LogOut", MiddleWare.Auth(), Controller.LogOut)
		account.POST("/Register", Controller.Register)
	}

	stuff := router.Group("/stuff")
	{
		stuff.GET("/getStuff", MiddleWare.Auth(), Controller.GetStuff)
	}
	router.Run(":8080")
}
