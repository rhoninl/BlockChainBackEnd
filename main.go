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
		account.GET("/Info", MiddleWare.Auth(), Controller.Info)
	}

	stuff := router.Group("/stuff")
	{
		stuff.GET("/getStuff", MiddleWare.Auth(), Controller.GetStuff)
	}
	company := router.Group("/company")
	{
		company.GET("/getJointVenture", MiddleWare.Auth(), Controller.GetJointVenture)
	}
	order := router.Group("/order")
	{
		order.GET("/getAllOrder", MiddleWare.Auth(), Controller.GetAllOrder)
	}
	router.Run(":8080")
}
