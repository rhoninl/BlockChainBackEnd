package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main/Controller"
	"main/MiddleWare"
)

const f = `
		⣿ ⣿ ⣿ ⣿ ⣿ ⠟ ⠋ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⢁ ⠈ ⢻ ⢿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⠃ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠈ ⡀ ⠭ ⢿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⡟ ⠄ ⢀ ⣾ ⣿ ⣿ ⣿ ⣷ ⣶ ⣿ ⣷ ⣶ ⣶ ⡆ ⠄ ⠄ ⠄ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⡇ ⢀ ⣼ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣧ ⠄ ⠄ ⢸ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣇ ⣼ ⣿ ⣿ ⠿ ⠶ ⠙ ⣿ ⡟ ⠡ ⣴ ⣿ ⣽ ⣿ ⣧ ⠄ ⢸ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣾ ⣿ ⣿ ⣟ ⣭ ⣾ ⣿ ⣷ ⣶ ⣶ ⣴ ⣶ ⣿ ⣿ ⢄ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⡟ ⣩ ⣿ ⣿ ⣿ ⡏ ⢻ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣹ ⡋ ⠘ ⠷ ⣦ ⣀ ⣠ ⡶ ⠁ ⠈ ⠁ ⠄ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣍ ⠃ ⣴ ⣶ ⡔ ⠒ ⠄ ⣠ ⢀ ⠄ ⠄ ⠄ ⡨ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣿ ⣦ ⡘ ⠿ ⣷ ⣿ ⠿ ⠟ ⠃ ⠄ ⠄ ⣠ ⡇ ⠈ ⠻ ⣿ ⣿ ⣿ ⣿
		⣿ ⣿ ⣿ ⣿ ⡿ ⠟ ⠋ ⢁ ⣷ ⣠ ⠄ ⠄ ⠄ ⠄ ⣀ ⣠ ⣾ ⡟ ⠄ ⠄ ⠄ ⠄ ⠉ ⠙ ⠻
		⡿ ⠟ ⠋ ⠁ ⠄ ⠄ ⠄ ⢸ ⣿ ⣿ ⡯ ⢓ ⣴ ⣾ ⣿ ⣿ ⡟ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄
		⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⣿ ⡟ ⣷ ⠄ ⠹ ⣿ ⣿ ⣿ ⡿ ⠁ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄ ⠄
`

func main() {
	fmt.Println(f)
	router := gin.Default()
	router.Use(cors.Default())
	homePage := router.Group("/homePage")
	{
		homePage.GET("/index", Controller.GetIndex)
	}
	account := router.Group("/account")
	{
		account.POST("/AuthCode", Controller.GetAuth)
		account.POST("/Login", Controller.Login)
		account.POST("/Register", Controller.Register)
		account.POST("/LogOut", MiddleWare.Auth(), Controller.LogOut)
		account.GET("/Info", MiddleWare.Auth(), Controller.Info)
		account.POST("/EditInfo", MiddleWare.Auth(), Controller.EditInfo)
		account.POST("/Forget", Controller.ForgetPassword)
		account.POST("/ChangePd", MiddleWare.Auth(), Controller.ChangePassword)
	}

	stuff := router.Group("/stuff")
	{
		stuff.GET("/getStuff", MiddleWare.Auth(), Controller.GetStuff)
		stuff.POST("/addStuff", MiddleWare.Auth(), Controller.AddStuff)
	}
	company := router.Group("/company")
	{
		company.GET("/getJointVenture", MiddleWare.Auth(), Controller.GetJointVenture)
		company.POST("/makeFriends", MiddleWare.Auth(), Controller.MakeFriend)
	}
	order := router.Group("/order")
	{
		order.GET("/getAllOrder", MiddleWare.Auth(), Controller.GetAllOrder)
		order.POST("/submitOrder", MiddleWare.Auth(), Controller.BindForm)
	}
	message := router.Group("/message")
	{
		message.GET("/getMessage", MiddleWare.Auth(), Controller.GetMessage)
		message.GET("/getMessageInfo/:id", MiddleWare.Auth(), Controller.GetMessageInfo)
		message.GET("/deleteMessage/:messageId", MiddleWare.Auth(), Controller.DeleteMessage)
	}
	router.GET("/ws", MiddleWare.Auth(), Controller.BuildSocket)
	router.Run(":8080")
}
